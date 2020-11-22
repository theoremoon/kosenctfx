<template>
  <div>
    <div>
      <p>
        <input type="file" @change="loadChallengeFiles" webkitdirectory />
        <input
          type="submit"
          value="Add Challenges"
          v-bind:disabled="challengeFiles.length == 0"
          @click="addChallenges"
        />
      </p>
    </div>

    <h3 class="text-xl">Loaded Challenges</h3>
    <ul class="ml-4">
      <li v-for="n in challengeFiles" :key="n">{{ n }}</li>
    </ul>

    <h3 class="text-xl">Message Log</h3>
    <div class="ml-4 message-list">
      <p v-for="(m, i) in message_log" :key="i">{{ m }}</p>
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
import { errorHandle } from "@/message";
import md5 from "blueimp-md5";
import { message } from "../../../message";

export default Vue.extend({
  data() {
    return { tree: null, challengeFiles: [], message_log: [] };
  },
  methods: {
    // Parseまわりの挙動
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
    async uploadFile(filename, data) {
      return API.post("admin/get-presigned-url", {
        key: filename
      })
        .then(r => {
          return fetch(r.data.presignedURL, {
            method: "PUT",
            body: new Blob([data])
          }).then(() => {
            return r.data.downloadURL;
          });
        })
        .catch(e => {
          errorHandle(this, e);
        });
    },
    async addChallenges() {
      await this.addChallenges_impl();
      this.$emit("update");
    },
    async addChallenges_impl() {
      const add_promises = [];

      for (let path of this.challengeFiles) {
        const taskFile = this.getFile(this.tree, path);
        const taskText = await readAsText(taskFile.file);
        const taskInfo = yaml.safeLoad(taskText);
        taskInfo["description"] = format(taskInfo["description"], taskInfo);
        taskInfo["attachments"] = [];

        const dirPath = path
          .split("/")
          .slice(0, -1)
          .join("/");
        const distfilesPath = dirPath + "/distfiles";
        const distfiles = this.getFile(this.tree, distfilesPath);
        if (distfiles) {
          const file_promises = [];
          this.travarseFileTree(distfiles, f => {
            file_promises.push(
              readAsText(f).then(fileData => ({
                header: {
                  name: f.webkitRelativePath.slice(distfilesPath.length + 1),
                  mtime: new Date(0)
                },
                buffer: fileData
              }))
            );
            return true;
          });

          const files = await Promise.all(file_promises);
          files.sort((a, b) => a.header.name.localeCompare(b.header.name));

          const pack = tar.pack();
          const pack_promises = [];
          files.forEach(f => {
            this.message_log.push("Read file: " + f.header.name);
            pack.entry(f.header, f.buffer);
          });

          await Promise.all(pack_promises);
          pack.finalize();

          const tarBlob = await streamToBlob(pack);
          const tarData = await tarBlob.arrayBuffer();
          const Datenow = Date.now;
          // DIRTY
          Date.now = () => {
            return 0;
          };
          const tarGZData = new Zlib.Gzip(new Uint8Array(tarData), {
            flags: { fname: undefined, comment: false }
          }).compress();
          Date.now = Datenow;

          const dataDigest = md5(tarGZData);
          const filename = format(
            "{0}_{1}.tar.gz",
            dirPath.split("/").slice("-1")[0],
            dataDigest
          );
          const url = await this.uploadFile(filename, tarGZData);
          taskInfo["attachments"].push({
            name: filename,
            url: url
          });
          this.message_log.push("Uploaded attachment: " + filename);
        }

        const rawDistfilesPath = dirPath + "/rawdistfiles";
        const rawDistfiles = this.getFile(this.tree, rawDistfilesPath);
        if (rawDistfiles) {
          const promises = [];
          this.travarseFileTree(rawDistfiles, f => {
            promises.push(
              readAsText(f).then(fileData => {
                const dataDigest = md5(fileData);
                const rawFilename = f.webkitRelativePath
                  .split("/")
                  .slice(-1)[0];
                const filename = format("{0}_{1}", rawFilename, dataDigest);
                this.uploadFile(filename, fileData).then(url => {
                  taskInfo["attachments"].push({
                    name: filename,
                    url: url
                  });
                  this.message_log.push("Uploaded attachment: " + filename);
                });
              })
            );
          });
        }

        add_promises.push(
          API.post("admin/new-challenge", taskInfo)
            .then(r => {
              message(this, r.data);
            })
            .catch(e => errorHandle(this, e))
        );
      }
      return Promise.all(add_promises);
    }
  }
});
</script>

<style lang="scss" scoped>
.message-list {
  height: 10rem;
  overflow-y: scroll;
}
</style>
