package atheon.scanners;

import atheon.core.*;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.*;

public class OpenAiScanner implements Scanner {
    private static final Pattern PATTERN =
        Pattern.compile("\\bsk-[A-Za-z0-9_-]{20,}\\b");

    public String name() {
        return "openai-api-key";
    }

    public String description() {
        return "Detects OpenAI API keys";
    }

    public Severity severity() {
        return Severity.CRITICAL;
    }

    public List<String> scan(String input) {
        List<String> matches = new ArrayList<>();
        Matcher m = PATTERN.matcher(input);

        while (m.find()) {
            matches.add(m.group());
        }

        return matches;
    }
}
