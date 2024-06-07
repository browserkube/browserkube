package io.browserkube;

import org.junit.jupiter.api.*;
import org.openqa.selenium.By;
import org.openqa.selenium.bidi.LogInspector;
import org.openqa.selenium.bidi.log.ConsoleLogEntry;
import org.openqa.selenium.chrome.ChromeOptions;
import org.openqa.selenium.remote.RemoteWebDriver;
import org.openqa.selenium.remote.SessionId;

import java.io.IOException;
import java.net.URI;
import java.net.URL;
import java.util.*;

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
public class WebDriverBiDiTests {
    int testCounter = 1;
    Properties props;
    private static final String SELENIUM_CUSTOM_PAGE = "https://www.selenium.dev/selenium/web/bidi/logEntryAdded.html";
    private List<ConsoleLogEntry> consoleLogEntries;

    RemoteWebDriver driver;

    @BeforeAll
    void init() throws IOException {
        System.setProperty("webdriver.chrome.driver", "drivers/chromedriver-mac-64bit");
        this.props = WebDriverExampleTests.readProps();
    }
    @BeforeEach
    void testInit() throws IOException {
        this.consoleLogEntries = new ArrayList<>();

        Map<String, Object> browserkubeOptions = new HashMap<>();
        browserkubeOptions.put("enableVNC", true);
        browserkubeOptions.put("tenant", "test");
        browserkubeOptions.put("name", "java test session " + this.testCounter);
        ChromeOptions chromeOptions = new ChromeOptions();
        //BiDi flag
        chromeOptions.setCapability("webSocketUrl", true);
        //
        chromeOptions.setCapability("browserName", "chrome");

        chromeOptions.setCapability("browserkube:options", browserkubeOptions);

         driver = (RemoteWebDriver) RemoteWebDriver.builder()
                .address(new URL(props.getProperty("WD_URL")))
                .addAlternative(chromeOptions)
                .build();
    }

    @AfterEach()
    void testTeardown() {
        this.consoleLogEntries.clear();
        if (this.driver != null) {
            driver.quit();
        }
        this.testCounter++;
    }

    @Test
    public void testBiDiWSConn() throws Exception {
        SessionId sessionID = driver.getSessionId();
        System.out.println("SessionID :" + sessionID.toString());
        //After the session is opened connect to the bidi websocket server

        String bidiURL = "ws://" + props.getProperty("BROWSERKUBE_URL") + "/browserkube/wd/hub/bidi/" + sessionID;
        System.out.println("WS URL:" + bidiURL);
        WebSocketClient c = new WebSocketClient(new URI(bidiURL));
        c.connect();
        for (int i = 0; i < 3; i++) {
            if (c.isOpen())
                break;
            Thread.sleep(250);
        }
        // Make sure that bidi channel is open
        c.send("Totally valid command");
        boolean isMsgReceived = false;
        for (int i = 0; i < 3; i++) {
            if(c.getMessage().isEmpty()) {
                Thread.sleep(250);
                continue;
            }
            isMsgReceived = true;
        }
        Assertions.assertTrue(isMsgReceived);
    }
    @Test
    public void testBiDiLogInspection() throws Exception {
        LogInspector logInspector = new LogInspector(driver);
        logInspector.onConsoleLog(log -> consoleLogEntries.add(log));

        driver.get(SELENIUM_CUSTOM_PAGE);

        driver.findElement(By.id("consoleLog")).click();

        // Wait until the log comes through the BiDi comm
        for (int i = 0; i < 3; i++) {
            if (!this.consoleLogEntries.isEmpty()) {
                break;
            }
            Thread.sleep(250);
        }

        Assertions.assertNotEquals(0, consoleLogEntries.size());
        ConsoleLogEntry consoleLogEntry = consoleLogEntries.get(0);
        Assertions.assertEquals("Hello, world!", consoleLogEntry.getText());
        Assertions.assertNull(consoleLogEntry.getRealm());
        Assertions.assertEquals("log", consoleLogEntry.getMethod());

        driver.quit();
    }
}
