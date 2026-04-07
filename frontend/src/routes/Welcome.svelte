<script lang="ts">
  import { push } from '@keenmate/svelte-spa-router';
  import { ValidateApiKey } from '../../wailsjs/go/main/App.js';
  import { shadowConfig, templateCache } from '../lib/store';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import EyeIcon from '@lucide/svelte/icons/eye';
  import EyeOffIcon from '@lucide/svelte/icons/eye-off';
  import { BrowserOpenURL } from '../../wailsjs/runtime/runtime.js';

  let buttonLabel = $state('Verify');
  let displayBadApiKeyErr = $state(false);
  let displayConnectionErr = $state(false);
  let displayValidationError = $state(false);
  let buttonActive = $state(true);
  let accessToken = $state(($shadowConfig['AccessToken'] as string) ?? '');
  let showToken = $state(false);

  function openURL() {
    BrowserOpenURL('https://app.safetyculture.com/account/api-tokens');
  }

  function validate() {
    if (accessToken.length === 0) {
      displayValidationError = true;
      return;
    }

    buttonActive = false;
    ValidateApiKey(accessToken)
      .then((result: string) => {
        if (result !== '') {
          buttonLabel = 'Try again';
          if (result === 'connection error') {
            displayConnectionErr = true;
          } else {
            displayBadApiKeyErr = true;
          }
          displayValidationError = true;
        } else {
          displayValidationError = false;
          if ($shadowConfig['AccessToken'] !== accessToken) {
            ($shadowConfig as any)['Export'] = ($shadowConfig as any)['Export'] ?? {};
            ($shadowConfig as any)['Export']['TemplateIds'] = [];
          }
          templateCache.set([]);
          $shadowConfig['AccessToken'] = accessToken;
          push('/config');
        }
      })
      .finally(() => {
        buttonActive = true;
      });
  }
</script>

<div class="flex h-full flex-col items-center">
  <img class="w-[120px] pt-8" src="../images/logo.svg" alt="SafetyCulture logo" />
  <h1 class="mt-4 mb-2 text-xl font-semibold">Welcome to SafetyCulture Exporter</h1>
  <img class="w-[500px]" src="../images/welcome.png" alt="welcome" />
  <p class="mt-4 mb-3 text-sm text-muted-foreground">
    Generate an API token from your
    <button class="cursor-pointer text-primary hover:underline" onclick={openURL}
      >SafetyCulture account</button
    >.
  </p>

  <div class="flex items-start gap-2">
    <div class="relative w-[340px]">
      <Input
        type={showToken ? 'text' : 'password'}
        placeholder="Enter API Token here"
        bind:value={accessToken}
        class="pr-9 {displayValidationError
          ? 'border-destructive focus-visible:ring-destructive/50'
          : ''}"
      />
      <button
        type="button"
        class="absolute top-1/2 right-3 -translate-y-1/2 cursor-pointer text-muted-foreground hover:text-foreground"
        onclick={() => (showToken = !showToken)}
      >
        {#if showToken}
          <EyeOffIcon class="size-4" />
        {:else}
          <EyeIcon class="size-4" />
        {/if}
      </button>
    </div>
    <Button onclick={validate} disabled={!buttonActive}>
      {buttonLabel}
    </Button>
  </div>

  {#if displayBadApiKeyErr}
    <div class="mt-2 max-w-md text-xs">
      <p class="font-medium text-destructive">Unable to verify token</p>
      <p class="text-muted-foreground">
        It looks like your token may be expired after 30 days of inactivity. Please generate a new
        token and try again.
      </p>
    </div>
  {/if}

  {#if displayConnectionErr}
    <div class="mt-2 max-w-md text-xs">
      <p class="font-medium text-destructive">Connection error</p>
      <p class="text-muted-foreground">
        It looks like you are not connected to the internet or behind a firewall. Please check your
        connection and try again.
      </p>
    </div>
  {/if}

  <section class="mt-4 max-w-lg">
    <div
      class="flex gap-3 rounded-lg border border-border bg-muted/50 p-3 text-xs text-muted-foreground"
    >
      <img src="../images/warning.svg" alt="alert icon" class="size-5 shrink-0" />
      <div>
        <p class="pb-1 text-xs font-semibold text-foreground">Important note</p>
        <p class="leading-5">
          All files (apart from SQL) you export will be stored in the same place on your computer or
          server as the SafetyCulture Exporter. If you want to change where your files get exported,
          please move the SafetyCulture Exporter file itself to that place.
        </p>
      </div>
    </div>
  </section>
</div>
