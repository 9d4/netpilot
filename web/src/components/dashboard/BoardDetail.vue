<script lang="ts" setup>
import { useBoardsApi } from '@/composables/api'
import type { Board } from '@/types/board'
import { AxiosError } from 'axios'
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const open = ref(false)
const fetching = ref(true)
const updating = ref(false)
const board = ref({} as Board)
const showPwd = ref(false)
const disabled = ref(true)
const updateErrors = ref<any>({})
const updateSuccess = ref(false)
const connCheckStatus = ref('')

const props = defineProps({
  uuid: { type: String }
})
const pwdToggle = () => {
  if (disabled.value) return
  showPwd.value = !showPwd.value
}

const editToggle = () => {
  disabled.value = !disabled.value
}

const close = () => {
  open.value = false
  setTimeout(() => {
    router.push({ params: { uuid: undefined } })
  }, 100)
}

const fetchBoard = () => {
  return useBoardsApi({ params: { detail: 1 } })
    .getOne(props.uuid!)
    .then((res) => {
      board.value = res.data
      board.value.password = atob(board.value.password)
    })
}

const updateBoard = () => {
  updateErrors.value = {}
  updateSuccess.value = false
  updating.value = true
  disabled.value = true
  useBoardsApi({
    data: board.value,
    headers: {
      'Content-Type': 'application/json'
    }
  })
    .put(props.uuid!)
    .then(() => {
      updateSuccess.value = true
      fetchBoard()
    })
    .catch((err: AxiosError) => {
      if (err.code === AxiosError.ERR_BAD_REQUEST) {
        updateErrors.value = err.response?.data!
      }
    })
    .finally(() => {
      updating.value = false
      disabled.value = false
    })
}

const checkConnection = () => {
  connCheckStatus.value = 'loading...'
  useBoardsApi()
    .extend('check', {
      method: 'post',
      data: {
        ...board.value
      }
    })
    .then(({ status }) => {
      if (status === 200) {
        connCheckStatus.value = 'OK'
      }
    })
    .catch((err: AxiosError) => {
      let msg = (m: any) => `Fail: ${m}`
      if (err.code != AxiosError.ERR_NETWORK) {
        connCheckStatus.value = msg(err.response?.status)
        return
      }

      connCheckStatus.value = msg(err.message)
    })
}

onMounted(() => {
  setTimeout(() => {
    open.value = true
  }, 100)

  fetchBoard().finally(() => {
    fetching.value = false
  })
})
</script>

<template>
  <div class="modal" :class="{ 'modal-open': open }">
    <div class="modal-box w-11/12 max-w-5xl z-20">
      <button class="btn btn-sm btn-circle absolute right-4 top-4" @click="close">✕</button>
      <h3 class="font-bold text-lg">Board</h3>
      <p class="text-sm text-gray-500">{{ props.uuid }}</p>

      <p class="text-gray-300 my-4" v-if="fetching">Fething data...</p>
      <div class="mt-10" v-else>
        <button class="flex ml-auto btn btn-sm" @click="editToggle">
          <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" class="bi bi-pencil mr-2"
            viewBox="0 0 16 16">
            <path
              d="M12.146.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1 0 .708l-10 10a.5.5 0 0 1-.168.11l-5 2a.5.5 0 0 1-.65-.65l2-5a.5.5 0 0 1 .11-.168l10-10zM11.207 2.5 13.5 4.793 14.793 3.5 12.5 1.207 11.207 2.5zm1.586 3L10.5 3.207 4 9.707V10h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.293l6.5-6.5zm-9.761 5.175-.106.106-1.528 3.821 3.821-1.528.106-.106A.5.5 0 0 1 5 12.5V12h-.5a.5.5 0 0 1-.5-.5V11h-.5a.5.5 0 0 1-.468-.325z" />
          </svg>
          Edit
        </button>

        <div class="form-control w-full my-4">
          <label class="label">
            <span class="label-text">Board Name</span>
          </label>
          <input :disabled="disabled" v-model="board.name" type="text" class="input input-sm input-bordered w-full" />
        </div>

        <div class="form-control w-full my-4 grid grid-cols-2 gap-4">
          <div>
            <label class="label">
              <span class="label-text">Host</span>
            </label>
            <input :disabled="disabled" v-model="board.host" type="text" class="input input-sm input-bordered w-full" />
          </div>
          <div>
            <label class="label">
              <span class="label-text">REST Port</span>
            </label>
            <input :disabled="disabled" v-model="board.port" type="text" class="input input-sm input-bordered w-full" />
          </div>
        </div>

        <div class="form-control w-max">
          <label class="label cursor-pointer gap-2">
            <input :disabled="disabled" type="checkbox" v-model="board.insecure_skip_verify"
              class="checkbox checkbox-sm" />
            <span class="label-text">Skip Insecure Verification</span>
          </label>
        </div>

        <div class="form-control w-full my-4 grid grid-cols-2 gap-4">
          <div>
            <label class="label">
              <span class="label-text">Username</span>
            </label>
            <input :disabled="disabled" v-model="board.user" type="text" class="input input-sm input-bordered w-full" />
          </div>
          <div>
            <label class="label">
              <span class="label-text">Password</span>
            </label>
            <div class="input-group">
              <input :disabled="disabled" v-model="board.password" :type="showPwd ? 'text' : 'password'"
                class="input input-sm input-bordered w-full" />
              <button class="btn btn-sm btn-square" @click="pwdToggle">
                <svg v-if="showPwd" xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                  class="bi bi-eye" viewBox="0 0 16 16">
                  <path
                    d="M16 8s-3-5.5-8-5.5S0 8 0 8s3 5.5 8 5.5S16 8 16 8zM1.173 8a13.133 13.133 0 0 1 1.66-2.043C4.12 4.668 5.88 3.5 8 3.5c2.12 0 3.879 1.168 5.168 2.457A13.133 13.133 0 0 1 14.828 8c-.058.087-.122.183-.195.288-.335.48-.83 1.12-1.465 1.755C11.879 11.332 10.119 12.5 8 12.5c-2.12 0-3.879-1.168-5.168-2.457A13.134 13.134 0 0 1 1.172 8z" />
                  <path d="M8 5.5a2.5 2.5 0 1 0 0 5 2.5 2.5 0 0 0 0-5zM4.5 8a3.5 3.5 0 1 1 7 0 3.5 3.5 0 0 1-7 0z" />
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                  class="bi bi-eye-slash" viewBox="0 0 16 16">
                  <path
                    d="M13.359 11.238C15.06 9.72 16 8 16 8s-3-5.5-8-5.5a7.028 7.028 0 0 0-2.79.588l.77.771A5.944 5.944 0 0 1 8 3.5c2.12 0 3.879 1.168 5.168 2.457A13.134 13.134 0 0 1 14.828 8c-.058.087-.122.183-.195.288-.335.48-.83 1.12-1.465 1.755-.165.165-.337.328-.517.486l.708.709z" />
                  <path
                    d="M11.297 9.176a3.5 3.5 0 0 0-4.474-4.474l.823.823a2.5 2.5 0 0 1 2.829 2.829l.822.822zm-2.943 1.299.822.822a3.5 3.5 0 0 1-4.474-4.474l.823.823a2.5 2.5 0 0 0 2.829 2.829z" />
                  <path
                    d="M3.35 5.47c-.18.16-.353.322-.518.487A13.134 13.134 0 0 0 1.172 8l.195.288c.335.48.83 1.12 1.465 1.755C4.121 11.332 5.881 12.5 8 12.5c.716 0 1.39-.133 2.02-.36l.77.772A7.029 7.029 0 0 1 8 13.5C3 13.5 0 8 0 8s.939-1.721 2.641-3.238l.708.709zm10.296 8.884-12-12 .708-.708 12 12-.708.708z" />
                </svg>
              </button>
            </div>
          </div>
        </div>

        <div class="form-control w-full my-4">
          <div class="input-group">
            <input v-model="connCheckStatus" disabled type="text" class="input input-sm input-bordered w-full" />
            <button class="btn btn-sm" @click="checkConnection">Check Connection</button>
          </div>
        </div>

        <div class="my-4">
          <div v-if="updateErrors.errors != undefined" class="alert alert-error my-2">
            <ul class="block">
              <li v-for="a in updateErrors?.errors" :key="a.field">{{ a.field }} is {{ a.tag }}</li>
            </ul>
          </div>
          <div class="flex justify-end items-center gap-5">
            <span v-if="updateSuccess" class="text-gray-500 text-sm">Successful</span>
            <button class="btn btn-sm" @click="updateBoard" v-if="!disabled">
              {{ updating ? 'Updating' : 'Update' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
