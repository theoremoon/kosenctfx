<template>
  <div>
    <div class="flex justify-around">
      <form class="lg:w-1/4 py-2 text-lg">
        <div class="flex">
          <input
            type="text"
            placeholder="KosenCTF{[\x21-\x7e]+}"
            class="w-full text-center focus:outline-none flaginput"
          />
          <Button class="flex-shrink-0 ml-1">Submit</Button>
        </div>
      </form>

      <div class="lg:w-1/4 py-2 text-lg">
        filter: <input type="text" class="text-center focus:outline-none flaginput" v-model="filter">
      </div>
    </div>

    <div class="flex flex-wrap">
      <template v-for="c in challenges">
        <template v-if="!filter || c.name.includes(filter) || tagfilter(c.tags, filter)">
        <!-- in detail (focused) mode -->
        <div
          class="mx-2 lg:w-1/3 challenge"
          :key="c.name"
          @click="focus = null"
          v-if="focus == c.name"
        >
          <div class="challenge-name">
            <span class="challenge-name-bg font-bold">
              <font-awesome-icon icon="flag" v-if="c.solvedby.includes(myteam)" />
              {{c.name}}
            </span>
          </div>
          <div class="challenge-info p-2 flex">
            <div class="flex-grow">
              <div>
                <span>{{c.score}}pts/{{c.solvedby.length}}solves</span>
              </div>
              <div>
                <span v-for="t in c.tags" :key="t" class="tag">{{t}}</span>
              </div>

              <div class="my-2">
                <div v-html="c.description"></div>
                <div class="text-right">author:{{c.author}}</div>
              </div>

              <div v-if="c.attachments && c.attachments.length > 0">
                <span v-for="a in c.attachments" :key="a.name">
                  <a :href="a.url" download class="rounded py-2 px-4 dl-button" @click.stop>
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
                  <li v-for="t in c.solvedby" :key="t" class="solvedteam">{{t}}</li>
                </ul>
              </div>
            </div>
          </div>
        </div>

        <!-- not focused (normal) -->
        <div class="mx-2 lg:w-1/5 challenge" :key="c.name" @click="focus = c.name" v-else>
          <div class="text-center challenge-name">
            <span class="challenge-name-bg font-bold">
              <font-awesome-icon icon="flag"  v-if="c.solvedby.includes(myteam)" />
              {{c.name}}
            </span>
          </div>
          <div class="challenge-info p-2">
            <div class="text-center">{{c.score}}pts/{{c.solvedby.length}}solves</div>
            <div class="text-center">
              <span v-for="t in c.tags" :key="t" class="tag">{{t}}</span>
            </div>
          </div>
        </div>
        </template>
      </template>
    </div>
  </div>
</template>

<script>
import Button from "@/components/Button";
export default {
  components: {
    Button,
  },
  data() {
    return {
      filter: "",
      myteam: "zer0pts",
      focus: null,
      challenges: [
        {
          name: "padrsa",
          description: "I added padding to the RSA.",
          score: 500,
          solvedby: [
            "zer0pts",
            "Harekaze",
            "TokyoWesterns",
            "TSG",
            "shibad0gs",
            "StarrySky",
            "binja",
            "noraneko"
          ],
          attachments: [
            {
              url:
                "https://file-examples.com/wp-content/uploads/2017/02/file-sample_100kB.doc",
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
          solvedby: ["zer0pts", "Harekaze", "soooooooooolong team name"],
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
  },
  methods: {
    tagfilter(tags, filter) {
      for (const t of tags) {
        if (t.includes(filter)) {
          return true;
        }
      }
      return false;
    }
  },
};
</script>

<style lang="scss" scoped>
@import "@/vars.scss";

.tag {
  display: inline-block;

  font-weight: bold;
  border-radius: 0.25rem;
  background-color: #718096;
  padding: 0 0.25rem;
  margin-right: 0.25rem;
}
.challenge {
  cursor: pointer;
}
.challenge-name {
  line-height: normal;

  .challenge-name-bg {
    display: inline-block;
    max-width: 90%;

    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;

    position: relative;
    background-color: $bg-color;
    z-index: 1000;
  }

  &:before {
    content: "";
    display: block;
    border: 2px solid $fg-color;
    border-bottom: none;

    height: 1rem;

    position: relative;
    bottom: -1.75rem;

    z-index: 500;
  }
}
.challenge-info {
  border: 2px solid $fg-color;
  border-top: none;
}
.solvedby {
  overflow-y: auto;
  .solvedby-list {
    height: 0;
  }
}

.flaginput {
  background: transparent;
  border-bottom: 1px solid $light-color;
}

.solvedteam {
  border-bottom: 1px solid $fg-color;
}

.dl-button {
  border: 1px solid $light-color;
  &:hover {
    background-color: $light-color;
  }
}
</style>