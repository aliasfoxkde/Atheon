package atheon.cli;

import atheon.core.*;
import atheon.output.Printer;
import picocli.CommandLine;
import picocli.CommandLine.*;

import java.nio.file.*;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.Callable;

@Command(name = "atheon", mixinStandardHelpOptions = true, version = "atheon 1.0",
         description = "Secret and credential scanner",
         subcommands = {Main.ScanCommand.class, Main.ListCommand.class})
public class Main implements Callable<Integer> {

    public static void main(String[] args) {
        System.exit(new CommandLine(new Main()).execute(args));
    }

    @Override
    public Integer call() {
        CommandLine.usage(this, System.out);
        return 0;
    }

    @Command(name = "scan", description = "Scan a file, directory, env vars, or stdin")
    static class ScanCommand implements Callable<Integer> {

        @Parameters(index = "0", description = "Path to scan, --env, or --stdin", defaultValue = "")
        String target;

        @Option(names = "--json", description = "Output findings as JSON")
        boolean json;

        @Option(names = "--exclude", description = "Comma-separated directory names to skip", split = ",")
        List<String> exclude = new ArrayList<>();

        @Option(names = "--ext", description = "Comma-separated extensions to scan", split = ",")
        List<String> extensions = new ArrayList<>();

        @Override
        public Integer call() throws Exception {
            if (target.isEmpty()) {
                System.err.println("atheon: scan requires a path, --env, or --stdin");
                return 1;
            }

            Registry registry = new Registry();
            Runner runner = new Runner(registry);
            runner.setExclude(exclude);
            runner.setExtensions(normalizeExts(extensions));

            List<Finding> findings;

            if (target.equals("--env")) {
                findings = runner.scanEnv();
            } else if (target.equals("--stdin")) {
                findings = runner.scanReader(System.in);
            } else {
                Path path = Paths.get(target);
                if (!Files.exists(path)) {
                    System.err.println("atheon: path not found: " + target);
                    return 1;
                }
                findings = Files.isDirectory(path)
                    ? runner.scanDir(path)
                    : runner.scanFile(path);
            }

            if (json) {
                Printer.printJson(findings);
            } else {
                Printer.print(findings, runner.getStats());
            }

            return findings.isEmpty() ? 0 : 1;
        }

        private List<String> normalizeExts(List<String> exts) {
            List<String> out = new ArrayList<>();
            for (String e : exts) out.add(e.startsWith(".") ? e : "." + e);
            return out;
        }
    }

    @Command(name = "list", description = "List all registered scanners")
    static class ListCommand implements Callable<Integer> {

        @Override
        public Integer call() {
            Registry registry = new Registry();
            List<Scanner> all = registry.all();
            System.out.printf("registered scanners (%d)%n%n", all.size());
            System.out.printf("  %-30s  %s%n", "NAME", "DESCRIPTION");
            System.out.printf("  %-30s  %s%n", "─".repeat(30), "─".repeat(43));
            for (Scanner s : all) {
                System.out.printf("  %-30s  %s%n", s.name(), s.description());
            }
            return 0;
        }
    }
}
