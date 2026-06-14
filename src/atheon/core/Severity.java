package atheon.core;

public enum Severity {
    CRITICAL("critical"),
    HIGH("high"),
    MEDIUM("medium"),
    LOW("low");

    private final String value;

    Severity(String value) { this.value = value; }

    @Override
    public String toString() { return value; }
}
