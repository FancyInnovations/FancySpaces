package com.fancyinnovations.fancyspaces.versions;

import com.google.gson.annotations.SerializedName;

import java.time.Instant;
import java.util.Date;
import java.util.List;

public record Version(
        @SerializedName("space_id") String spaceID,
        String id,
        String name,
        String platform,
        String channel,
        @SerializedName("published_at") String publishedAt,
        String changelog,
        @SerializedName("supported_platform_versions") List<String> supportedPlatformVersions,
        List<VersionFile> files
) {

    public long publishedAtMillis() {
        return Instant.parse(publishedAt).toEpochMilli();
    }

    public Date publishedAtDate() {
        return new Date(publishedAtMillis());
    }

    @Override
    public String toString() {
        return "Version{" +
                "spaceID='" + spaceID + '\'' +
                ", id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", platform='" + platform + '\'' +
                ", channel='" + channel + '\'' +
                ", publishedAt='" + publishedAt + '\'' +
                ", changelog='" + changelog + '\'' +
                ", supportedPlatformVersions=" + supportedPlatformVersions +
                ", files=" + files +
                '}';
    }
}
