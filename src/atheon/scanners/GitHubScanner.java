package atheon.scanners;

import atheon.core.*;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.*;

public class GitHubScanner implements Scanner {
    private static final Pattern PATTERN = Pattern.compile("ghp_[0-9a-zA-Z]{36}");

    public String name()        { return "github-pat"; }
    public String description() { return "Detects GitHub personal access tokens (ghp_...)"; }
    public Severity severity()  { return Severity.CRITICAL; }

    public List<String> scan(String input) {
        List<String> matches = new ArrayList<>();
        Matcher m = PATTERN.matcher(input);
        while (m.find()) matches.add(m.group());
        return matches;
    }
}
