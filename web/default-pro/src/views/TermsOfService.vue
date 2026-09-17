<template>
  <LegalLayout
    :title="$t('legal.termsTitle')"
    :description="$t('legal.termsDescription')"
    :updatedAt="$t('legal.updatedAt')"
  >
    <div class="legal-intro">
      <p>{{ $t('legal.termsIntro1') }}</p>
      <p>{{ $t('legal.termsIntro2') }}</p>
    </div>

    <section class="legal-section">
      <h2>{{ $t('legal.termsS1Title') }}</h2>
      <p>{{ $t('legal.termsS1P1') }}</p>
      <p>{{ $t('legal.termsS1P2') }}</p>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.termsS2Title') }}</h2>
      <p>{{ $t('legal.termsS2P1') }}</p>
      <ul>
        <li>{{ $t('legal.termsS2L1') }}</li>
        <li>{{ $t('legal.termsS2L2') }}</li>
        <li>{{ $t('legal.termsS2L3') }}</li>
      </ul>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.termsS3Title') }}</h2>
      <p>{{ $t('legal.termsS3P1') }}</p>
      <ul>
        <li>{{ $t('legal.termsS3L1') }}</li>
        <li>{{ $t('legal.termsS3L2') }}</li>
        <li>{{ $t('legal.termsS3L3') }}</li>
        <li>{{ $t('legal.termsS3L4') }}</li>
      </ul>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.termsS4Title') }}</h2>
      <p>{{ $t('legal.termsS4P1') }}</p>
      <p>{{ $t('legal.termsS4P2') }}</p>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.termsS5Title') }}</h2>
      <p>{{ $t('legal.termsS5P1') }}</p>
      <p>{{ $t('legal.termsS5P2') }}</p>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.termsS6Title') }}</h2>
      <p>{{ $t('legal.termsS6P1') }}</p>
      <p>{{ $t('legal.termsS6P2') }}</p>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.termsS7Title') }}</h2>
      <p>{{ $t('legal.termsS7P1') }}</p>
      <p>{{ $t('legal.termsS7P2') }}</p>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.termsS8Title') }}</h2>
      <p>{{ $t('legal.termsS8P1') }}</p>
      <p>{{ $t('legal.termsS8P2') }}</p>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.termsS9Title') }}</h2>
      <p>{{ $t('legal.termsS9P1') }}</p>
      <p>{{ $t('legal.termsS9P2') }}</p>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.termsS10Title') }}</h2>
      <p>{{ $t('legal.termsS10P1') }}</p>
    </section>
  </LegalLayout>
</template>

<script setup>
import LegalLayout from '@/layouts/LegalLayout.vue'
</script>
