interface Board {
  created_at: string;
  updated_at: string;
  uuid: string;
  name: string;
  host: string;
  port: string;
  insecure_skip_verify: boolean;
  user: string;
}

export const useBoard = () => {
  const config = useRuntimeConfig();
  const { data } = useFetch<{ boards: Board[] }>('/boards', {
    baseURL: config.public.apiBase,
  });

  const boards = data.value?.boards;
  return { boards };
};
