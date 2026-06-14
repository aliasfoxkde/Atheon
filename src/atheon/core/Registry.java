package atheon.core;

import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import java.nio.file.*;
import java.util.*;
import java.util.stream.Stream;

public class Registry {

    private final List<Scanner> scanners;

    public Registry() {
        this.scanners = Collections.unmodifiableList(discover());
    }

    public List<Scanner> all() { return scanners; }
    public int count()         { return scanners.size(); }

    private List<Scanner> discover() {
        List<Scanner> found = new ArrayList<>();
        try {
            URI uri = Registry.class.getProtectionDomain().getCodeSource().getLocation().toURI();
            Path root;
            if (uri.getScheme().equals("file") && uri.getPath().endsWith(".jar")) {
                FileSystem fs = FileSystems.newFileSystem(URI.create("jar:" + uri), Map.of());
                root = fs.getPath("/");
            } else {
                root = Path.of(uri);
            }
            try (Stream<Path> stream = Files.walk(root)) {
                stream.filter(p -> p.toString().endsWith(".class"))
                      .forEach(p -> {
                          String className = toClassName(root, p);
                          if (className != null) tryLoad(className, found);
                      });
            }
        } catch (URISyntaxException | IOException e) {
            throw new RuntimeException("Failed to discover scanners", e);
        }
        found.sort(Comparator.comparing(Scanner::name));
        return found;
    }

    private String toClassName(Path root, Path classFile) {
        Path rel = root.relativize(classFile);
        String s = rel.toString().replace('\\', '/');
        if (!s.startsWith("atheon/scanners/")) return null;
        return s.replace('/', '.').replace(".class", "");
    }

    private void tryLoad(String className, List<Scanner> out) {
        try {
            Class<?> cls = Class.forName(className);
            if (cls.isInterface() || !Scanner.class.isAssignableFrom(cls)) return;
            out.add((Scanner) cls.getDeclaredConstructor().newInstance());
        } catch (Exception ignored) {}
    }
}
