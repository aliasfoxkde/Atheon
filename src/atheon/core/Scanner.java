package atheon.core;

import java.util.List;

public interface Scanner {
    String name();
    String description();
    Severity severity();
    List<String> scan(String input);
}
