package com.fancyinnovations.fancyspaces;

import com.fancyinnovations.fancyspaces.versions.VersionService;
import de.oliver.fancyanalytics.logger.ExtendedFancyLogger;

public class FancySpaces {

    private final ExtendedFancyLogger fancyLogger;

    private final String baseURL;

    private final VersionService versionService;

    public FancySpaces(String apiKey) {
        this.fancyLogger = new ExtendedFancyLogger("FancySpaces Java-SDK");

        this.baseURL = "https://fancyspaces.net/api/v1";

        this.versionService = new VersionService(this, apiKey == null ? "" : apiKey);
    }

    public FancySpaces() {
        this("");
    }

    public String getBaseURL() {
        return baseURL;
    }

    public ExtendedFancyLogger getFancyLogger() {
        return fancyLogger;
    }

    public VersionService getVersionService() {
        return versionService;
    }
}
