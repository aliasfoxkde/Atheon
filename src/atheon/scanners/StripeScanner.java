package atheon.scanners;

import atheon.core.*;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.*;

public class StripeScanner implements Scanner {
    private static final Pattern PATTERN = Pattern.compile("sk_live_[0-9a-zA-Z]{24}");

    public String name()        { return "stripe-secret-key"; }
    public String description() { return "Detects Stripe secret keys (sk_live_...)"; }
    public Severity severity()  { return Severity.CRITICAL; }

    public List<String> scan(String input) {
        List<String> matches = new ArrayList<>();
        Matcher m = PATTERN.matcher(input);
        while (m.find()) matches.add(m.group());
        return matches;
    }
}
