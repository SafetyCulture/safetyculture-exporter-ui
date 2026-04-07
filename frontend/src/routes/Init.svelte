<script lang="ts">
  import { push } from '@keenmate/svelte-spa-router';
  import { ReadVersion, GetSettings, ValidateApiKey } from '../../wailsjs/go/main/App.js';
  import { shadowConfig, latestVersion, emptyStores } from '../lib/store';
  import { isNullOrEmptyObject } from '../lib/utils';

  async function processVersion(): Promise<boolean> {
    return await ReadVersion()
      .then((result: Record<string, string> | null) => {
        if (result != null) {
          latestVersion.set(result);
        } else {
          latestVersion.set({});
        }
      })
      .then(() => {
        return shouldForceUpdate() === true;
      })
      .catch(() => {
        return false;
      });
  }

  async function processAccessToken(): Promise<string> {
    if ('AccessToken' in $shadowConfig) {
      const token = $shadowConfig['AccessToken'] as string;
      if (token.length === 0) {
        return '/welcome';
      } else {
        return ValidateApiKey(token).then((res: string) => {
          if (res !== '') {
            return '/welcome';
          } else {
            return '/config';
          }
        });
      }
    } else {
      return '/welcome';
    }
  }

  function shouldForceUpdate(): boolean {
    return (
      !isNullOrEmptyObject($latestVersion) &&
      $latestVersion['should_update'] === true &&
      $latestVersion['current'] !== 'v0.0.0-dev' &&
      $latestVersion['download_url'] !== ''
    );
  }

  emptyStores();

  GetSettings()
    .then((result: Record<string, unknown>) => {
      shadowConfig.set(result);
    })
    .then(() => {
      processVersion().then((shouldUpdate: boolean) => {
        if (shouldUpdate === true) {
          push('/update');
        } else {
          processAccessToken().then((page: string) => {
            push(page);
          });
        }
      });
    });
</script>
