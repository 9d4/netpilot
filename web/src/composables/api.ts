import axios from 'axios'
import type { AxiosRequestConfig } from 'axios'

interface Options extends AxiosRequestConfig {}

export const useApi = (resource: string, opts?: Options) => {
  const baseUrl = import.meta.env.VITE_API_BASE_URL || `${location.origin}/api`
  const url =`${baseUrl}/${resource}`

  return {
    get: () => {
      return axios(url, opts)
    },
    post: () => {
      return axios(url, {
        method: 'post',
        ...opts
      })
    },
    put: (uuid:string) => {
      return axios(`${url}/${uuid}`, {
        method: 'put',
        ...opts
      })
    },
    delete: (uuid: string) => {
      return axios(`${url}/${uuid}`, {
        method: 'delete',
        ...opts
      })
    }
  }
}
