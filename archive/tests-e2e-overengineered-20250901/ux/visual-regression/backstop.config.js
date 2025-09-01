const SEPTEO_COLORS = {
  primary: '#ff6136',
  primaryDark: '#e55a2e',
  secondary: '#2c3e50',
  success: '#27ae60',
  warning: '#f39c12',
  danger: '#e74c3c'
};

module.exports = {
  "id": "fire-salamander-visual-tests",
  "paths": {
    "bitmaps_reference": "ux/visual-regression/backstop_data/bitmaps_reference",
    "bitmaps_test": "ux/visual-regression/backstop_data/bitmaps_test",
    "engine_scripts": "ux/visual-regression/backstop_data/engine_scripts",
    "html_report": "ux/visual-regression/backstop_data/html_report",
    "ci_report": "ux/visual-regression/backstop_data/ci_report"
  },
  "viewports": [
    {
      "label": "phone",
      "width": 375,
      "height": 667
    },
    {
      "label": "tablet",
      "width": 768,
      "height": 1024
    },
    {
      "label": "desktop",
      "width": 1920,
      "height": 1080
    }
  ],
  "onBeforeScript": "puppet/onBefore.js",
  "onReadyScript": "puppet/onReady.js",
  "scenarios": [
    {
      "label": "Homepage - Dashboard",
      "cookiePath": "backstop_data/engine_scripts/cookies.json",
      "url": "http://localhost:8080",
      "referenceUrl": "",
      "readyEvent": "",
      "readySelector": ".main-content",
      "delay": 2000,
      "hideSelectors": [],
      "removeSelectors": [],
      "hoverSelector": "",
      "clickSelector": "",
      "postInteractionWait": 0,
      "selectors": [
        "document",
        ".navbar",
        ".main-content",
        ".stats-grid",
        ".charts-section"
      ],
      "misMatchThreshold": 0.1,
      "requireSameDimensions": true
    },
    {
      "label": "Analyzer Page",
      "url": "http://localhost:8080",
      "readySelector": "#analyzer-page",
      "clickSelector": "a[data-page='analyzer']",
      "postInteractionWait": 1000,
      "selectors": [
        "#analyzer-page",
        ".analyzer-form",
        ".form-container"
      ],
      "misMatchThreshold": 0.1
    },
    {
      "label": "Analysis Results",
      "url": "http://localhost:8080",
      "readySelector": "#analyzer-page",
      "clickSelector": "a[data-page='analyzer']",
      "postInteractionWait": 1000,
      "keyPressSelectors": [
        {
          "selector": "#analysisUrl",
          "keyPress": "https://example.com"
        }
      ],
      "clickSelectors": [
        "#analysisForm button[type='submit']"
      ],
      "postInteractionWait": 5000,
      "selectors": [
        "#analysisResults",
        ".results-container"
      ],
      "misMatchThreshold": 0.2
    },
    {
      "label": "History Page",
      "url": "http://localhost:8080",
      "clickSelector": "a[data-page='history']",
      "postInteractionWait": 1000,
      "selectors": [
        "#history-page",
        ".history-content",
        ".search-box"
      ],
      "misMatchThreshold": 0.1
    },
    {
      "label": "Reports Page",
      "url": "http://localhost:8080",
      "clickSelector": "a[data-page='reports']",
      "postInteractionWait": 1000,
      "selectors": [
        "#reports-page",
        ".report-generator",
        ".saved-reports"
      ],
      "misMatchThreshold": 0.1
    },
    {
      "label": "Monitoring Page",
      "url": "http://localhost:8080",
      "clickSelector": "a[data-page='monitoring']",
      "postInteractionWait": 1000,
      "selectors": [
        "#monitoring-page",
        ".health-section",
        ".metrics-section"
      ],
      "misMatchThreshold": 0.1
    },
    {
      "label": "SEPTEO Color Validation",
      "url": "http://localhost:8080",
      "onReadyScript": "puppet/validateSepteoColors.js",
      "selectors": [
        ".navbar",
        ".btn-primary",
        ".status-indicator"
      ],
      "misMatchThreshold": 0.05
    },
    {
      "label": "Responsive Navigation - Mobile",
      "url": "http://localhost:8080",
      "viewports": [{"label": "phone", "width": 375, "height": 667}],
      "selectors": [
        ".navbar",
        ".nav-menu"
      ],
      "misMatchThreshold": 0.1
    },
    {
      "label": "Dark Mode Compatibility",
      "url": "http://localhost:8080",
      "onReadyScript": "puppet/enableDarkMode.js",
      "selectors": [
        "document"
      ],
      "misMatchThreshold": 0.2
    },
    {
      "label": "Loading States",
      "url": "http://localhost:8080",
      "clickSelector": "a[data-page='analyzer']",
      "keyPressSelectors": [
        {
          "selector": "#analysisUrl",
          "keyPress": "https://example.com"
        }
      ],
      "clickSelectors": [
        "#analysisForm button[type='submit']"
      ],
      "postInteractionWait": 1000,
      "selectors": [
        "#analysisProgress",
        ".progress-container"
      ],
      "misMatchThreshold": 0.3
    }
  ],
  "paths": {
    "bitmaps_reference": "backstop_data/bitmaps_reference",
    "bitmaps_test": "backstop_data/bitmaps_test",
    "engine_scripts": "backstop_data/engine_scripts",
    "html_report": "backstop_data/html_report",
    "ci_report": "backstop_data/ci_report"
  },
  "report": ["browser"],
  "engine": "puppeteer",
  "engineOptions": {
    "args": ["--no-sandbox"]
  },
  "asyncCaptureLimit": 5,
  "asyncCompareLimit": 50,
  "debug": false,
  "debugWindow": false
};