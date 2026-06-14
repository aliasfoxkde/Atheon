package atheon.core;

public class Finding {
    public final String scanner;
    public final String severity;
    public final String file;
    public final Integer line;
    public final String description;
    public final String match;

    public Finding(String scanner, String severity, String file, Integer line, String description, String match) {
        this.scanner = scanner;
        this.severity = severity;
        this.file = file;
        this.line = line;
        this.description = description;
        this.match = match;
    }

    public Finding withFile(String file) {
        return new Finding(scanner, severity, file, line, description, match);
    }
}
