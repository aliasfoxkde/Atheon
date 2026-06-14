package atheon.scanners;

import atheon.core.*;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.*;

public class TwilioScanner implements Scanner {
    private static final Pattern PATTERN = Pattern.compile("AC[0-9a-fA-F]{32}");

    public String name()        { return "twilio-account-sid"; }
    public String description() { return "Detects Twilio account SIDs (AC...)"; }
    public Severity severity()  { return Severity.HIGH; }

    public List<String> scan(String input) {
        List<String> matches = new ArrayList<>();
        Matcher m = PATTERN.matcher(input);
        while (m.find()) matches.add(m.group());
        return matches;
    }
}
