package atheon.core;

import java.io.*;
import java.nio.charset.StandardCharsets;
import java.nio.file.*;
import java.nio.file.attribute.BasicFileAttributes;
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicLong;

public class Runner {

    private final Registry registry;
    private List<String> exclude   = new ArrayList<>();
    private List<String> extensions = new ArrayList<>();
    private Stats stats = new Stats();

    private static final Set<String> SKIP_DIRS = Set.of(
        ".git", "node_modules", "vendor", ".terraform", "dist", "build", "__pycache__"
    );

    private static final Set<String> BINARY_EXTS = Set.of(
        ".png", ".jpg", ".jpeg", ".gif", ".pdf", ".zip", ".tar", ".gz",
        ".exe", ".bin", ".so", ".dylib", ".class", ".jar"
    );

    public Runner(Registry registry) { this.registry = registry; }

    public void setExclude(List<String> exclude)       { this.exclude = exclude; }
    public void setExtensions(List<String> extensions) { this.extensions = extensions; }
    public Stats getStats()                            { return stats; }

    public List<Finding> scanString(String content) {
        return scanLines(content);
    }

    public List<Finding> scanFile(Path path) throws IOException {
        long start = System.currentTimeMillis();
        byte[] data = Files.readAllBytes(path);
        List<Finding> findings = new ArrayList<>();
        for (Finding f : scanLines(new String(data, StandardCharsets.UTF_8))) {
            findings.add(f.withFile(path.toString()));
        }
        stats = new Stats();
        stats.files    = 1;
        stats.bytes    = data.length;
        stats.elapsedMs = System.currentTimeMillis() - start;
        return findings;
    }

    public List<Finding> scanDir(Path root) throws IOException {
        long start = System.currentTimeMillis();
        List<Path> paths = new ArrayList<>();

        Files.walkFileTree(root, new SimpleFileVisitor<>() {
            @Override
            public FileVisitResult preVisitDirectory(Path dir, BasicFileAttributes attrs) {
                String name = dir.getFileName() != null ? dir.getFileName().toString() : "";
                if (SKIP_DIRS.contains(name) || exclude.contains(name))
                    return FileVisitResult.SKIP_SUBTREE;
                return FileVisitResult.CONTINUE;
            }

            @Override
            public FileVisitResult visitFile(Path file, BasicFileAttributes attrs) {
                if (!skipFile(file)) paths.add(file);
                return FileVisitResult.CONTINUE;
            }
        });

        int threads = Math.max(1, Runtime.getRuntime().availableProcessors());
        ExecutorService pool = Executors.newFixedThreadPool(threads);
        List<Future<List<Finding>>> futures = new ArrayList<>();
        AtomicLong totalBytes = new AtomicLong();

        for (Path p : paths) {
            futures.add(pool.submit(() -> {
                try {
                    byte[] data = Files.readAllBytes(p);
                    totalBytes.addAndGet(data.length);
                    List<Finding> found = new ArrayList<>();
                    for (Finding f : scanLines(new String(data, StandardCharsets.UTF_8))) {
                        found.add(f.withFile(p.toString()));
                    }
                    return found;
                } catch (IOException e) {
                    return List.of();
                }
            }));
        }

        pool.shutdown();
        try {
            pool.awaitTermination(5, TimeUnit.MINUTES);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }

        List<Finding> findings = new ArrayList<>();
        for (Future<List<Finding>> future : futures) {
            try { findings.addAll(future.get()); } catch (Exception ignored) {}
        }

        stats = new Stats();
        stats.files    = paths.size();
        stats.bytes    = totalBytes.get();
        stats.elapsedMs = System.currentTimeMillis() - start;
        return findings;
    }

    public List<Finding> scanEnv() {
        List<Finding> findings = new ArrayList<>();
        for (Map.Entry<String, String> entry : System.getenv().entrySet()) {
            for (Scanner s : registry.all()) {
                for (String match : s.scan(entry.getValue())) {
                    findings.add(new Finding(s.name(), s.severity().toString(), "env:" + entry.getKey(), null, s.description(), match));
                }
            }
        }
        return findings;
    }

    public List<Finding> scanReader(InputStream in) throws IOException {
        byte[] data = in.readAllBytes();
        List<Finding> findings = new ArrayList<>();
        for (Finding f : scanLines(new String(data, StandardCharsets.UTF_8))) {
            findings.add(f.withFile("stdin"));
        }
        return findings;
    }

    private List<Finding> scanLines(String content) {
        List<Finding> findings = new ArrayList<>();
        String[] lines = content.split("\n", -1);
        for (int i = 0; i < lines.length; i++) {
            int lineNum = i + 1;
            for (Scanner s : registry.all()) {
                for (String match : s.scan(lines[i])) {
                    findings.add(new Finding(s.name(), s.severity().toString(), null, lineNum, s.description(), match));
                }
            }
        }
        return findings;
    }

    private boolean skipFile(Path path) {
        String name = path.toString();
        int dot = name.lastIndexOf('.');
        String ext = dot >= 0 ? name.substring(dot).toLowerCase() : "";
        if (BINARY_EXTS.contains(ext)) return true;
        if (extensions.isEmpty()) return false;
        return !extensions.contains(ext);
    }
}
