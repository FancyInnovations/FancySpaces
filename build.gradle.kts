plugins {
    id("com.gradleup.shadow") version "9.3.1" apply false
}

allprojects {
    repositories {
        mavenLocal()
        mavenCentral()
        maven("https://repo.papermc.io/repository/maven-public/")
        maven("https://repo.fancyinnovations.com/releases")
        maven(url = "https://jitpack.io")
    }
}
