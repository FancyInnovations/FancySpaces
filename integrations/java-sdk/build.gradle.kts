plugins {
    id("java-library")
    id("maven-publish")
}

allprojects {
    group = "com.fancyinnovations"
    version = getSDKVersion()
    description = "SDK for the FancySpaces api"

    repositories {
        mavenLocal()
        mavenCentral()
        maven(url = "https://repo.fancyinnovations.com/releases")
    }
}

dependencies {
    compileOnly("de.oliver.FancyAnalytics:logger:0.0.9")

    compileOnly("com.google.code.gson:gson:2.13.2")
    implementation("org.jetbrains:annotations:26.1.0")
}

tasks {
    publishing {
        repositories {
            maven {
                name = "fancyspacesReleases"
                url = uri("https://maven.fancyspaces.net/fancyinnovations/releases")

                credentials(HttpHeaderCredentials::class) {
                    name = "Authorization"
                    value = providers
                        .gradleProperty("fancyspacesApiKey")
                        .orElse(
                            providers
                                .environmentVariable("FANCYSPACES_API_KEY")
                                .orElse("")
                        )
                        .get()
                }

                authentication {
                    create<HttpHeaderAuthentication>("header")
                }
            }

            maven {
                name = "fancyspacesSnapshots"
                url = uri("https://maven.fancyspaces.net/fancyinnovations/snapshots")

                credentials(HttpHeaderCredentials::class) {
                    name = "Authorization"
                    value = providers
                        .gradleProperty("fancyspacesApiKey")
                        .orElse(
                            providers
                                .environmentVariable("FANCYSPACES_API_KEY")
                                .orElse("")
                        )
                        .get()
                }

                authentication {
                    create<HttpHeaderAuthentication>("header")
                }
            }

//            maven {
//                name = "fancyinnovationsReleases"
//                url = uri("https://repo.fancyinnovations.com/releases")
//                credentials(PasswordCredentials::class)
//                authentication {
//                    isAllowInsecureProtocol = true
//                    create<BasicAuthentication>("basic")
//                }
//            }
//
//            maven {
//                name = "fancyinnovationsSnapshots"
//                url = uri("https://repo.fancyinnovations.com/snapshots")
//                credentials(PasswordCredentials::class)
//                authentication {
//                    isAllowInsecureProtocol = true
//                    create<BasicAuthentication>("basic")
//                }
//            }
        }
        publications {
            create<MavenPublication>("maven") {
                groupId = "com.fancyinnovations.fancyspaces"
                artifactId = "java-sdk"
                version = getSDKVersion()
                from(project.components["java"])
            }
        }
    }

    compileJava {
        options.encoding = Charsets.UTF_8.name()
        options.release = 25
    }

    test {
        useJUnitPlatform()
    }

}

java {
    toolchain.languageVersion.set(JavaLanguageVersion.of(25))
    withJavadocJar()
    withSourcesJar()
}

fun getSDKVersion(): String {
    return file("VERSION").readText()
}
