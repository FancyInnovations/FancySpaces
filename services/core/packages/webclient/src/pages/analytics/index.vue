<script lang="ts" setup>
</script>

<template>
    <header class="hero">
        <v-container class="hero-container">
            <v-row align="center" justify="center">
                <v-col class="text-center" cols="12" md="8">
                    <v-img class="mb-8" height="200" src="https://fancyinnovations.com/logos-and-banners/fancyanalytics-banner.png"/>

                    <h1 class="text-h3 font-weight-medium">Explore your data like never before</h1>

                    <p class="text-subtitle-1 mb-8">
                      Collect, analyze, and visualize your data with ease.
                    </p>

                    <div class="hero-ctas">
                        <v-btn class="mr-4" color="primary" to="register">
                          Get started
                        </v-btn>

                        <v-btn color="primary" href="https://fancyinnovations.com/docs/web-services/fancyanalytics" variant="tonal">
                          View docs
                        </v-btn>
                    </div>
                </v-col>
            </v-row>
        </v-container>
    </header>

    <main>
        <v-container class="features-container mt-12 w-100">
            <v-row class="mb-4" justify="center">
                <v-col class="text-center mb-8" cols="4">
                    <h2 class="text-h4 font-weight-medium">Use cases</h2>
                    <p class="body 1">FancyAnalytics is built from the ground up to support all kinds of
                        applications, from Minecraft servers to GoLang apps. We have set a focus on Minecraft
                        and Hytale related applications tho.</p>
                </v-col>
            </v-row>

            <v-row class="mb-8" justify="center">
                <v-col class="mb-6" cols="12" md="2">
                    <v-card variant="text">
                        <v-card-text class="d-flex flex-column align-center text-center">
                            <v-icon class="mb-4" color="primary" size="48">mdi-server</v-icon>
                            <h3 class="mb-1">Minecraft Servers</h3>
                            <p>With the FancyAnalytics Minecraft plugin, you track player
                                related metrics or server performance easily.</p>
                        </v-card-text>
                    </v-card>
                </v-col>

                <v-col class="mb-6" cols="12" md="2">
                    <v-card variant="text">
                        <v-card-text class="d-flex flex-column align-center text-center">
                            <v-icon class="mb-4" color="primary" size="48">mdi-power-plug</v-icon>
                            <h3 class="mb-1">Minecraft Plugins</h3>
                            <p>With our SDKs, you can integrate FancyAnalytics into your
                                Minecraft plugins and track usage of your plugins.</p>
                        </v-card-text>
                    </v-card>
                </v-col>

                <v-col class="mb-6" cols="12" md="2">
                    <v-card variant="text">
                        <v-card-text class="d-flex flex-column align-center text-center">
                            <v-icon class="mb-4" color="primary" size="48">mdi-server</v-icon>
                            <h3 class="mb-1">Hytale Servers</h3>
                            <p>Host Hytale servers? Use FancyAnalytics to track player
                                metrics and server performance easily.</p>
                        </v-card-text>
                    </v-card>
                </v-col>

                <v-col class="mb-6" cols="12" md="2">
                    <v-card variant="text">
                        <v-card-text class="d-flex flex-column align-center text-center">
                            <v-icon class="mb-4" color="primary" size="48">mdi-power-plug</v-icon>
                            <h3 class="mb-1">Hytale Plugins</h3>
                            <p>Integrate FancyAnalytics into your Hytale plugins
                                using our SDKs and track usage of your plugins.</p>
                        </v-card-text>
                    </v-card>
                </v-col>
            </v-row>

            <v-row class="section-bg mb-8 py-8" justify="center">
                <v-col class="text-center mb-4 pt-8" cols="12">
                    <h2 class="text-h4 font-weight-medium">Features that empower your team</h2>
                    <p class="text-subtitle-1">From data ingestion to visualization, FancyAnalytics has you
                        covered.</p>
                </v-col>

                <v-col class="mb-6" cols="12" md="2">
                    <v-card variant="text">
                        <v-card-text class="d-flex flex-column align-center text-center">
                            <v-icon color="primary" size="48">mdi-chart-areaspline</v-icon>
                            <h3>Beautiful Visualizations</h3>
                            <p>Build dashboards with elegant charts and share insights in
                                seconds.</p>
                        </v-card-text>
                    </v-card>
                </v-col>

                <v-col class="mb-6" cols="12" md="2">
                    <v-card variant="text">
                        <v-card-text class="d-flex flex-column align-center text-center">
                            <v-icon color="primary" size="48">mdi-database</v-icon>
                            <h3>Easy Data Ingestion</h3>
                            <p>Use various SDKs and connectors to bring your data into
                                FancyAnalytics effortlessly.</p>
                        </v-card-text>
                    </v-card>
                </v-col>

                <v-col class="mb-6" cols="12" md="2">
                    <v-card variant="text">
                        <v-card-text class="d-flex flex-column align-center text-center">
                            <v-icon color="primary" size="48">mdi-update</v-icon>
                            <h3>Continuous Updates</h3>
                            <p>We will keep improving FancyAnalytics with new features and
                                enhancements.</p>
                        </v-card-text>
                    </v-card>
                </v-col>
            </v-row>

            <v-row class="mb-8" justify="center">
                <v-col cols="12" md="5">
                    <v-card class="callout" variant="tonal">
                        <v-card-text class="d-flex flex-row align-center justify-space-between">
                            <div>
                                <div class="callout-title">Ready to see FancyAnalytics in action?</div>
                                <div class="callout-sub">Sign up now and start exploring your data like never
                                    before.
                                </div>
                            </div>
                            <div>
                                <v-btn color="primary" to="register">create account</v-btn>
                            </div>
                        </v-card-text>
                    </v-card>
                </v-col>
            </v-row>
        </v-container>
    </main>
</template>

<style scoped>
.hero {
    padding: 100px 0;
    background: linear-gradient(135deg, rgba(241, 120, 99, 0.12) 0%, rgba(185, 146, 16, 0.06) 100%);
    backdrop-filter: blur(6px);
}

.features-container {
    min-width: 100%;
}

.section-bg {
    background-color: rgba(243, 173, 42, 0.18);
    border-radius: 8px;
}

</style>
