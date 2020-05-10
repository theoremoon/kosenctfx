<template>
  <div class="flex flex-wrap">
    <template v-for="c in challenges">
      <!-- in detail (focused) mode -->
      <div class="mx-2 lg:w-1/3" :key="c.name" @click="focus = null" v-if="focus == c.name">
        <div class="challenge-name">
          <span class="challenge-name-bg font-bold">
            <font-awesome-icon icon="flag" v-if="c.solvedby.includes(myteam)"/>
            {{c.name}}
            </span>
        </div>
        <div class="challenge-info p-2 flex">
          <div class="flex-grow">
            <div>
              <span>{{c.score}}pts/{{c.solvedby.length}}solves</span>
              <span v-for="t in c.tags" :key="t">[{{t}}]</span>
            </div>

            <div class="my-2">
              <div>{{c.description|raw}}</div>
              <div class="text-right">author:{{c.author}}</div>
            </div>

            <div v-if="c.attachments && c.attachments.length > 0">
              <span v-for="a in c.attachments" :key="a.name">
                <a :href="a.url" download class="bg-gray-300 hover:bg-gray-500 rounded py-2 px-4">
                  <font-awesome-icon icon="download" />
                  {{a.name}}
                </a>
              </span>
            </div>
          </div>
            <div class="flex-grow-0 ml-2 solvedby">
              <div class="font-bold">Solved by</div>
              <div class="solvedby-list">
                <ul>
                  <li v-for="t in c.solvedby" :key="t" style="border-bottom: 1px solid #ccc;">{{t}}</li>
                </ul>
              </div>
            </div>
          </div>
        </div>

      <!-- not focused (normal) -->
      <div class="mx-2 lg:w-1/5" :key="c.name" @click="focus = c.name" v-else>
        <div class="text-center challenge-name">
          <span class="challenge-name-bg font-bold">
            <font-awesome-icon icon="flag" v-if="c.solvedby.includes(myteam)"/>
            {{c.name}}
          </span>
        </div>
        <div class="challenge-info p-2">
          <div class="text-center">{{c.score}}pts/{{c.solvedby.length}}solves</div>
          <div class="text-center">
            <span v-for="t in c.tags" :key="t">[{{t}}]</span>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script>
export default {
  data() {
    return {
      myteam: "zer0pts",
      focus: null,
      challenges: [
        {
          name: "padrsa",
          description: "I added padding to the RSA.",
          score: 500,
          solvedby: ["zer0pts", "Harekaze", "TokyoWesterns", "TSG", "shibad0gs", "StarrySky", "binja", "noraneko"],
          attachments: [
            {
              url: "https://example.com/",
              name: "distfiles.tar.gz"
            }
          ],
          tags: ["crypto", "warmup"],
          author: "theoremoon"
        },
        {
          name: "can you find me in the secret wowowowowo",
          description: "I added padding to the RSA.",
          score: 500,
          solvedby: [],
          attachments: [
            {
              url: "https://example.com/",
              name: "distfiles.tar.gz"
            }
          ],
          tags: ["crypto"],
          author: "theoremoon"
        },
        {
          name: "padrsa2",
          description: "I added padding to the RSA.",
          score: 500,
          solvedby: ["zer0pts", "Harekaze", "TokyoWesterns"],
          tags: ["crypto"],
          author: "theoremoon"
        },
        {
          name: "padrsa3",
          description: "I added padding to the RSA.",
          score: 500,
          solvedby: ["zer0pts", "Harekaze", "TokyoWesterns"],
          tags: ["crypto"],
          author: "theoremoon"
        },
        {
          name: "padrsa4",
          description: "I added padding to the RSA.",
          score: 500,
          solvedby: ["zer0pts", "Harekaze", "TokyoWesterns"],
          tags: ["crypto"],
          author: "theoremoon"
        },
        {
          name: "padrsa5",
          description: "I added padding to the RSA.",
          score: 500,
          solvedby: ["zer0pts", "Harekaze", "TokyoWesterns"],
          tags: ["crypto"],
          author: "theoremoon"
        },
        {
          name: "padrsa6",
          description: "I added padding to the RSA.",
          score: 500,
          solvedby: ["zer0pts", "Harekaze", "TokyoWesterns"],
          tags: ["crypto"],
          author: "theoremoon"
        },
        {
          name: "padrsa7",
          description: "I added padding to the RSA.",
          score: 500,
          solvedby: ["zer0pts", "Harekaze", "TokyoWesterns"],
          tags: ["crypto"],
          author: "theoremoon"
        }
      ]
    };
  }
};
</script>

<style lang="scss" scoped>
.challenge-name {
  line-height: normal;

  .challenge-name-bg {
    display: inline-block;
    max-width: 90%;

    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;

    position: relative;
    background-color: white;
    z-index: 1000;
  }

  &:before {
    content: "";
    display: block;
    border: 2px solid black;
    border-bottom: none;

    height: 1rem;

    position: relative;
    bottom: -1.75rem;

    z-index: 500;
  }
}
.challenge-info {
  border: 2px solid black;
  border-top: none;
}
.solvedby {
  overflow-y: scroll;
  .solvedby-list {
    height: 0;
  }
}
</style>