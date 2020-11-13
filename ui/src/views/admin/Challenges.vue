<template>
  <div class="my-4 mx-8">
    <h2 class="text-2xl">Challenges</h2>
    <div>
      Load Challenge Directory
      <input type="file" @change="loadChallengeFiles" webkitdirectory />
    </div>
    <div>
      <input
        type="submit"
        value="Add Challenges"
        v-bind:disabled="challengeFiles.length == 0"
        @click="addChallenges"
      />
    </div>
    <div>
      Loaded Challenges
      <ul>
        <li v-for="n in challengeFiles" :key="n">{{ n }}</li>
      </ul>
    </div>
  </div>
</template>

<script>
import Vue from "vue";
import yaml from "js-yaml";
import format from "string-template";
import tar from "tar-stream";
import { readAsText } from "promise-file-reader";
import streamToBlob from "stream-to-blob";
import { Zlib } from "zlibjs/bin/gzip.min";
import API from "@/api";
import { errorHandle } from "../../message";
import md5 from "blueimp-md5";

export default Vue.extend({
  data() {
    return {
      tree: null,
      challengeFiles: []
    };
  },
  methods: {
    makeFileTree(files) {
      let tree = {
        type: "directory",
        children: {}
      };
      files.forEach(f => {
        const parts = f.webkitRelativePath.split("/");
        let top = tree;
        parts.forEach((p, i) => {
          if (!Object.prototype.hasOwnProperty.call(top.children, p)) {
            if (i == parts.length - 1) {
              top.children[p] = {
                type: "file",
                file: f
              };
            } else {
              top.children[p] = {
                type: "directory",
                children: {}
              };
            }
          }
          top = top.children[p];
        });
      });
      return tree;
    },
    travarseFileTree(tree, callback) {
      if (tree.type === "file") {
        if (!callback(tree.file)) {
          return false;
        }
        return true;
      }

      for (const [, value] of Object.entries(tree.children)) {
        if (!this.travarseFileTree(value, callback)) {
          return false;
        }
      }
      return true;
    },
    loadChallengeFiles(ev) {
      this.tree = this.makeFileTree(ev.target.files);
      this.challengeFiles = [];
      this.travarseFileTree(this.tree, f => {
        if (f.name === "task.yml") {
          this.challengeFiles.push(f.webkitRelativePath);
        }
        return true;
      });
    },
    getFile(tree, path) {
      const dirs = path.split("/");
      let top = tree;
      try {
        dirs.forEach(p => {
          top = top.children[p];
        });
        return top;
      } catch (e) {
        return null;
      }
    },

    async addChallenges() {
      for (let path of this.challengeFiles) {
        const taskFile = this.getFile(this.tree, path);
        const taskText = await readAsText(taskFile.file);
        const taskInfo = yaml.safeLoad(taskText);
        taskInfo["description"] = format(taskInfo["description"], taskInfo);

        const dirPath = path
          .split("/")
          .slice(0, -1)
          .join("/");
        const distfilesPath = dirPath + "/distfiles";
        const distfiles = this.getFile(this.tree, distfilesPath);
        taskInfo["attachments"] = [];
        if (distfiles) {
          const pack = tar.pack();

          const promises = [];
          this.travarseFileTree(distfiles, f => {
            promises.push(
              readAsText(f).then(fileData => {
                pack.entry(
                  {
                    name: f.webkitRelativePath.slice(distfilesPath.length + 1)
                  },
                  fileData
                );
              })
            );
          });

          await Promise.all(promises);
          pack.finalize();

          const tarBlob = await streamToBlob(pack);
          const tarData = await tarBlob.arrayBuffer();
          const tarGZData = new Zlib.Gzip(new Uint8Array(tarData)).compress();
          const dataDigest = md5(tarGZData);
          const filename = format(
            "{0}_{1}.tar.gz",
            dirPath.split("/").slice("-1")[0],
            dataDigest
          );

          const url = await API.post("admin/get-presigned-url", {
            key: filename
          })
            .then(r => {
              const formData = new FormData();
              formData.append("file", new Blob([tarGZData]));
              return fetch(r.data.presignedURL, {
                method: "PUT",
                headers: {
                  accept: "multipart/form-data"
                },
                body: formData
              }).then(() => {
                return r.data.downloadURL;
              });
            })
            .catch(e => {
              errorHandle(this, e);
            });
          console.log(url);
        }
      }
    }
  }
});
</script>
