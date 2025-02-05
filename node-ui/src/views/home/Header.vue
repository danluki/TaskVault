<script setup lang="ts">
import { SettingOne, Logout, SunOne, GithubOne, Me } from '@icon-park/vue-next';
import {useToast} from "@/components/ui/toast";
import {logoutUser} from "@/utils";
import {useRouter} from "vue-router";
import {ref} from "vue";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger
} from "@/components/ui/dropdown-menu";
import ThemeSwitcher from "@/components/ThemeSwitcher.vue";

const {toast} = useToast()
const router = useRouter()

function logout() {
  logoutUser()
  toast({
    description: 'Logout successfully.',
  });
  router.push("/login")
}

const dropdownOpen = ref(false);

const toggleDropdown = () => {
  dropdownOpen.value = !dropdownOpen.value;
};


</script>

<template>
  <header>
    <div class="flex h-full w-full justify-center items-center">
      <div class="w-[25%]">
        <router-link to="/" class="flex text-center items-center">
          <img class="w-[2.5rem] mx-[0.5rem]" src="../../assets/syncra.png" alt="" />
          <span class="font-bold font-mono text-2xl pl-[0.5rem]">Syncra</span>
        </router-link>
      </div>
      <div class="w-[58%] flex flex-col items-end">
<!--        <nav class="font-bold">-->
<!--          <ul class="flex space-x-4">-->
<!--            <li>-->
<!--              <router-link to="/" class="text-black">Home</router-link>-->
<!--            </li>-->
<!--            <li>-->
<!--              <router-link to="/docs" class="text-black">Document</router-link>-->
<!--            </li>-->
<!--            <li>-->
<!--              <router-link to="/posts" class="text-black">Post</router-link>-->
<!--            </li>-->
<!--          </ul>-->
<!--        </nav>-->
      </div>
      <div class="w-[17%] flex justify-end space-x-4 pr-[1rem]">
        <button class="flex items-center justify-center">
          <a href="https://github.com/danluki/TaskVault" target="_blank">
            <github-one theme="outline" size="18" :fill="['#333']" />
          </a>
        </button>
        <ThemeSwitcher/>
        <div class="relative">
          <DropdownMenu>
            <DropdownMenuTrigger as-child>
              <button class="flex items-center justify-center" @click="toggleDropdown">
                <me theme="two-tone" size="18" :fill="['#333', '#50e3c2']" />
              </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>

            </DropdownMenuContent>
          </DropdownMenu>
          <div v-if="dropdownOpen"  class="absolute right-0 bg-white shadow-lg mt-2 rounded-lg w-48 p-2">
            <div class="flex items-center text-center text-lg font-bold mt-[1rem] mb-[0.5rem] mx-4">
              <SunOne theme="two-tone" size="24" :fill="['#333', '#f8e71c']" />
              <span class="ml-2">Hi admin</span>
            </div>
            <ul>
              <li class="px-4 py-2 cursor-pointer hover:bg-gray-100 flex items-center space-x-2">
                <SettingOne /> <span>Settings</span>
              </li>
              <li @click="logout" class="px-4 py-2 cursor-pointer hover:bg-gray-100 flex items-center space-x-2">
                <Logout /> <span>Logout</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<style scoped>

</style>