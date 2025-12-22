const axios = require('axios');

class HelpNowAnalytics {
    /**
     * @param {string} baseUrl - The base URL of the HelpNow API (e.g., http://api.helpnow.com)
     */
    constructor(baseUrl) {
        this.baseUrl = baseUrl;
        this.token = null;

        // Try to verify if we are in a browser environment to auto-detect token from URL
        if (typeof window !== 'undefined') {
            const params = new URLSearchParams(window.location.search);
            this.token = params.get('t') || params.get('token') || params.get('id');
        }
    }

    /**
     * Set the tracking token manually.
     * @param {string} token 
     */
    identify(token) {
        this.token = token;
    }

    async _sendEvent(type, additionalData = {}) {
        if (!this.token) {
            console.warn('HelpNowAnalytics: No token set. Skipping event:', type);
            return;
        }

        const payload = {
            token: this.token,
            event_type: type,
            meta: {
                user_agent: (typeof navigator !== 'undefined') ? navigator.userAgent : 'node',
                url: (typeof window !== 'undefined') ? window.location.href : '',
                referrer: (typeof document !== 'undefined') ? document.referrer : '',
                ...additionalData
            }
        };

        try {
            await axios.post(`${this.baseUrl}/public/analytics`, payload);
        } catch (error) {
            console.error('HelpNowAnalytics: Failed to send event:', error);
        }
    }

    trackImpression() {
        return this._sendEvent('impression');
    }

    trackClick() {
        return this._sendEvent('click');
    }

    trackConversion(value = 0) {
        return this._sendEvent('conversion', { value });
    }

    /**
     * Track a custom event.
     * @param {string} eventName - The name of the event (e.g. 'purchased_item')
     * @param {object} properties - Additional metadata for the event
     */
    track(eventName, properties = {}) {
        return this._sendEvent(eventName, properties);
    }
}

module.exports = HelpNowAnalytics;
