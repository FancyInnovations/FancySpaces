<script lang="ts" setup>

</script>

<template>
<v-container>
  <v-row>
    <v-col>
      <h1 class="text-h3 mt-8 mb-4">Java SDK</h1>

      <h2 class="text-h4 mb-2">Include the SDK in your project</h2>

      <p class="text-body-1 mb-2">To include the Java SDK in your project, you need to add the following dependency to your build.gradle.kts or pom.xml file.</p>

      <p class="text-body-1"><strong>Gradle:</strong></p>
      <pre>
        <code>
repositories {
  maven("https://repo.fancyinnovations.com/releases")
}

dependencies {
  compileOnly("com.fancyinnovations.fancyspaces:java-sdk:VERSION")
}
        </code>
      </pre>

      <p class="text-body-1"><strong>Maven:</strong></p>
      <pre>
        <code>
&lt;repository>
  &lt;id>fancyinnovations-releases&lt;/id>
  &lt;name>FancyInnovations Repository&lt;/name>
  &lt;url>https://repo.fancyinnovations.com/releases&lt;/url>
&lt;/repository>

&lt;dependency>
  &lt;groupId>com.fancyinnovations.fancyspaces&lt;/groupId>
  &lt;artifactId>java-sdk&lt;/artifactId>
  &lt;version>VERSION&lt;/version>
  &lt;scope>provided&lt;/scope>
&lt;/dependency>
        </code>
      </pre>

      <p class="text-body-1">Replace `VERSION` with the version of the API you want to use. You can find the latest version on the download pages or in the Maven repository.</p>

      <h2 class="text-h4 mt-6 mb-2">Using the SDK</h2>

      <p class="text-body-1">You can then start using the SDK, by using the services provided in the FancySpaces class. Example: new FancySpaces().getVersionService().getLatestVersion("my-space");</p>
    </v-col>
  </v-row>
</v-container>
</template>

<style scoped>

</style>
