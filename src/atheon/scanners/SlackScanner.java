package atheon.scanners;

import atheon.core.*;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.*;

public class SlackScanner implements Scanner {
    private static final Pattern PATTERN = Pattern.compile("xoxb-[0-9]{11}-[0-9]{11}-[0-9a-zA-Z]{24}");

    public String name()        { return "slack-bot-token"; }
    public String description() { return "Detects Slack bot tokens (xoxb-...)"; }
    public Severity severity()  { return Severity.HIGH; }

    public List<String> scan(String input) {
        List<String> matches = new ArrayList<>();
        Matcher m = PATTERN.matcher(input);
        while (m.find()) matches.add(m.group());
        return matches;
    }
}
