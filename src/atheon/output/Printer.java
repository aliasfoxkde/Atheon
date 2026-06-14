package atheon.output;

import com.fasterxml.jackson.databind.ObjectMapper;
import atheon.core.Finding;
import atheon.core.Stats;

import java.util.List;

public class Printer {

    private static final String RESET  = "\033[0m";
    private static final String RED    = "\033[31m";
    private static final String YELLOW = "\033[33m";
    private static final String BLUE   = "\033[34m";
    private static final String WHITE  = "\033[37m";
    private static final String GREEN  = "\033[32m";

    public static void print(List<Finding> findings, Stats stats) {
        if (findings.isEmpty()) {
            System.out.println(GREEN + "✓ No secrets detected" + RESET);
        } else {
            System.out.println();
            for (Finding f : findings) {
                System.out.printf("%s[%s]%s %s%n", colorFor(f.severity), f.severity.toUpperCase(), RESET, f.scanner);
                System.out.printf("  file:  %s%n", location(f));
                System.out.printf("  desc:  %s%n", f.description);
                System.out.printf("  match: %s%n", redact(f.match));
                System.out.println();
            }
            System.out.println("─────────────────────────────");
            System.out.printf("found %d potential secret(s)%n", findings.size());
        }
        if (stats != null && stats.files > 0) {
            System.out.printf("%nfiles: %d  size: %s  time: %dms%n",
                stats.files, formatBytes(stats.bytes), stats.elapsedMs);
        }
    }

    public static void printJson(List<Finding> findings) throws Exception {
        System.out.println(new ObjectMapper().writerWithDefaultPrettyPrinter().writeValueAsString(findings));
    }

    private static String location(Finding f) {
        if (f.file == null) return "(none)";
        if (f.line != null) return f.file + ":" + f.line;
        return f.file;
    }

    private static String colorFor(String severity) {
        return switch (severity) {
            case "critical" -> RED;
            case "high"     -> YELLOW;
            case "medium"   -> BLUE;
            default         -> WHITE;
        };
    }

    private static String redact(String s) {
        if (s == null) return "";
        if (s.length() <= 8) return "*".repeat(s.length());
        return s.substring(0, 4) + "*".repeat(s.length() - 8) + s.substring(s.length() - 4);
    }

    private static String formatBytes(long b) {
        if (b >= 1L << 20) return String.format("%.1f MB", (double) b / (1L << 20));
        if (b >= 1L << 10) return String.format("%.1f KB", (double) b / (1L << 10));
        return b + " B";
    }
}
