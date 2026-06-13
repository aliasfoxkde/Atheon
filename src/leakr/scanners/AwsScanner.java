package leakr.scanners;

import leakr.core.*;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.*;

public class AwsScanner implements Scanner {
    private static final Pattern PATTERN = Pattern.compile("\\b(?:AKIA|ASIA)[0-9A-Z]{16}\\b");

    public String name()        { return "aws-access-key"; }
    public String description() { return "Detects AWS access key IDs (AKIA...)"; }
    public Severity severity()  { return Severity.CRITICAL; }

    public List<String> scan(String input) {
        List<String> matches = new ArrayList<>();
        Matcher m = PATTERN.matcher(input);
        while (m.find()) matches.add(m.group());
        return matches;
    }
}
