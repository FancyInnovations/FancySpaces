package com.fancyinnovations.fancyspaces.versions;

public record VersionFile(
        String name,
        String url,
        long size
) {


    @Override
    public String toString() {
        return "VersionFile{" +
                "name='" + name + '\'' +
                ", url='" + url + '\'' +
                ", size=" + size +
                '}';
    }
}
