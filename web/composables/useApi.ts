import { UseFetchOptions } from "nuxt/app";

interface options extends UseFetchOptions<any> {}

export default (path: string, opts?: options) => {
  const config = useRuntimeConfig();

  return useFetch(path, {
    baseURL: config.public.apiBase,
    ...opts,
  });
};
