const HelpNowAnalytics = require('./index');

const analytics = new HelpNowAnalytics('http://localhost:8080');

// Get token from command line
const token = process.argv[2];

if (!token) {
    console.error("Please provide a token as an argument");
    process.exit(1);
}

analytics.identify(token);

(async () => {
    console.log("Sending Impression...");
    await analytics.trackImpression();

    console.log("Sending Click...");
    await analytics.trackClick();

    console.log("Sending Conversion ($50.00)...");
    await analytics.trackConversion(50.00);

    console.log("Done.");
})();
