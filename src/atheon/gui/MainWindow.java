package atheon.gui;

import javax.swing.*;
import java.awt.*;

public class MainWindow extends JFrame {

    private final JTabbedPane tabs;

    public MainWindow() {
        setTitle("atheon");
        setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
        setSize(920, 600);
        setMinimumSize(new Dimension(700, 420));
        getContentPane().setBackground(AtheonColors.BG);

        tabs = new JTabbedPane(JTabbedPane.TOP);
        tabs.setBackground(AtheonColors.BG);
        tabs.setForeground(AtheonColors.GREEN);
        tabs.setFont(new Font(Font.MONOSPACED, Font.PLAIN, 12));
        tabs.setBorder(BorderFactory.createEmptyBorder());

        UIManager.put("TabbedPane.selected",              AtheonColors.BG);
        UIManager.put("TabbedPane.background",            AtheonColors.BG);
        UIManager.put("TabbedPane.foreground",            AtheonColors.GREEN);
        UIManager.put("TabbedPane.darkShadow",            AtheonColors.BG);
        UIManager.put("TabbedPane.shadow",                AtheonColors.DIM);
        UIManager.put("TabbedPane.highlight",             AtheonColors.BG);
        UIManager.put("TabbedPane.light",                 AtheonColors.BG);
        UIManager.put("TabbedPane.focus",                 AtheonColors.GREEN);
        UIManager.put("TabbedPane.contentBorderInsets",   new Insets(0, 0, 0, 0));
        UIManager.put("TabbedPane.tabInsets",             new Insets(4, 10, 4, 10));

        addTab();

        add(tabs);
        setLocationRelativeTo(null);
        setVisible(true);
    }

    public void addTab() {
        int n = tabs.getTabCount() + 1;
        TerminalPanel panel = new TerminalPanel(this);
        tabs.addTab("  session " + n + "  ", panel);
        tabs.setSelectedComponent(panel);
        SwingUtilities.invokeLater(panel::focusInput);
    }
}
