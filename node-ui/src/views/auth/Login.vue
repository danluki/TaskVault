<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card } from '@/components/ui/card'
import { ToastAction } from '@/components/ui/toast'
import { useToast } from '@/components/ui/toast/use-toast'

import { reactive, h } from 'vue'
import {useRouter} from "vue-router";
import {setUser, type User} from "@/utils"
import { User as UserIcon, RectangleEllipsis } from 'lucide-vue-next'

const router = useRouter()
const {toast} = useToast()
const user = reactive<User>({ Username: "", Password: "" })

const login = () => {
  if (user.Username === "admin" && user.Password === "admin") {
    setUser({
      Username: "admin",
      Password: "admin",
    })
    toast({
      description: 'Login successfully.',
    });
    router.push('/')
  } else {
    toast({
      title: 'Uh oh! Something went wrong.',
      description: 'Wrong username or password.',
      variant: 'destructive',
      action: h(ToastAction, {
        altText: 'Try again',
      }, {
        default: () => 'Try again',
      }),
    });
  }
}



</script>

<template>
  <div class="h-screen login-bg bg-slate-50">
    <div class="flex h-full justify-center items-center">
      <div class="h-max min-w-[16rem] w-1/4 max-w-[24rem] text-center">
        <div class="inline-flex mt-4 mb-8 items-center">
          <img src="../../../../docs/src/assets/syncra.png" class="h-12 mr-2" />
          <h1 class="font-bold text-4xl font-mono">Syncra</h1>
        </div>

        <Card class="p-6 shadow-lg">
          <form @submit.prevent="login">
            <div class="mb-3 relative w-full max-w-sm items-center">
              <Input id="user" v-model="user.Username" class="pl-10 w-full mt-1" placeholder="admin" />
              <span class="absolute start-0 inset-y-0 flex items-center justify-center px-2">
                <UserIcon class="size-6 text-muted-foreground" />
              </span>
            </div>

            <div class="mb-3 relative w-full max-w-sm items-center">
              <Input id="password" v-model="user.Password" type="password" class="pl-10 w-full mt-1" placeholder="" />
              <span class="absolute start-0 inset-y-0 flex items-center justify-center px-2">
                <RectangleEllipsis class="size-6 text-muted-foreground" />
              </span>
            </div>


            <Button type="submit" class="w-full mt-3">SIGN IN</Button>
          </form>
        </Card>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-bg {
  background-image: url('@/assets/login-bg.svg');
  background-repeat: no-repeat;
  background-size: 100% auto;
  background-position: 0 100%;
}
</style>