<template>
  <LegalLayout
    :title="$t('legal.privacyTitle')"
    :description="$t('legal.privacyDescription')"
    :updatedAt="$t('legal.updatedAt')"
  >
    <div class="legal-intro">
      <p>{{ $t('legal.privacyIntro1') }}</p>
      <p>{{ $t('legal.privacyIntro2') }}</p>
    </div>

    <section class="legal-section">
      <h2>{{ $t('legal.privacyS1Title') }}</h2>
      <p>{{ $t('legal.privacyS1P1') }}</p>
      <p>{{ $t('legal.privacyS1P2') }}</p>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.privacyS2Title') }}</h2>
      <p>{{ $t('legal.privacyS2P1') }}</p>
      <ul>
        <li>{{ $t('legal.privacyS2L1') }}</li>
        <li>{{ $t('legal.privacyS2L2') }}</li>
        <li>{{ $t('legal.privacyS2L3') }}</li>
        <li>{{ $t('legal.privacyS2L4') }}</li>
        <li>{{ $t('legal.privacyS2L5') }}</li>
      </ul>
      <p>{{ $t('legal.privacyS2P2') }}</p>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.privacyS3Title') }}</h2>
      <p>{{ $t('legal.privacyS3P1') }}</p>
      <ul>
        <li>{{ $t('legal.privacyS3L1') }}</li>
        <li>{{ $t('legal.privacyS3L2') }}</li>
        <li>{{ $t('legal.privacyS3L3') }}</li>
        <li>{{ $t('legal.privacyS3L4') }}</li>
        <li>{{ $t('legal.privacyS3L5') }}</li>
      </ul>
      <p>{{ $t('legal.privacyS3P2') }}</p>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.privacyS4Title') }}</h2>
      <p>{{ $t('legal.privacyS4P1') }}</p>
      <p>{{ $t('legal.privacyS4P2') }}</p>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.privacyS5Title') }}</h2>
      <p>{{ $t('legal.privacyS5P1') }}</p>
      <ul>
        <li>{{ $t('legal.privacyS5L1') }}</li>
        <li>{{ $t('legal.privacyS5L2') }}</li>
        <li>{{ $t('legal.privacyS5L3') }}</li>
        <li>{{ $t('legal.privacyS5L4') }}</li>
      </ul>
      <p>{{ $t('legal.privacyS5P2') }}</p>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.privacyS6Title') }}</h2>
      <p>{{ $t('legal.privacyS6P1') }}</p>
      <p>{{ $t('legal.privacyS6P2') }}</p>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.privacyS7Title') }}</h2>
      <p>{{ $t('legal.privacyS7P1') }}</p>
      <p>{{ $t('legal.privacyS7P2') }}</p>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.privacyS8Title') }}</h2>
      <p>{{ $t('legal.privacyS8P1') }}</p>
      <p>{{ $t('legal.privacyS8P2') }}</p>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.privacyS9Title') }}</h2>
      <p>{{ $t('legal.privacyS9P1') }}</p>
      <p>{{ $t('legal.privacyS9P2') }}</p>
    </section>

    <section class="legal-section">
      <h2>{{ $t('legal.privacyS10Title') }}</h2>
      <p>{{ $t('legal.privacyS10P1') }}</p>
    </section>
  </LegalLayout>
</template>

<script setup>
import LegalLayout from '@/layouts/LegalLayout.vue'
</script>
