<template>
  <div class="container-xl">
    <div class="row align-items-center">
      <div class="col-lg">
        <div class="form-floating">
          <textarea class="form-control" placeholder="Source code here" id="input" style="height: 400px" v-model="input"></textarea>
          <label for="input">Input</label>
        </div>
      </div>
      <div class="col-lg">
        <div class="form-floating">
          <textarea class="form-control" placeholder="Source code here" id="output" style="height: 400px" v-model="output"></textarea>
          <label for="output">Output</label>
        </div>
      </div>
      <div class="col-lg">
        <div class="form-floating">
          <textarea class="form-control" placeholder="Source code here" id="syntax-errors" style="height: 400px" v-model="syntax_errors"></textarea>
          <label for="syntax-errors">Syntax errors</label>
        </div>
      </div>
    </div>
    <div class="row">
      <div class="col-sm-2">
        <button type="button" class="btn btn-outline-primary" @click="ping">Ping server</button>
      </div>
      <div class="col-sm-2">
        <button type="button" class="btn btn-outline-primary" @click="checkForErrors">Check for errors</button>
      </div>
    </div>
  </div>
</template>

<script>
import server from '../services/server'


export default {
  name: 'GUI',
  props: {
    msg: String
  },
  data() {
    return {
      input: "",
      output: "",
      lexical_errors: "",
      syntax_errors: ""
    }
  },

  // methods
  methods: {
    async ping() {
      await server.ping().then(res => {
        this.output = res.data.message;
        console.log(this.output);
      });
    },
    async checkForErrors() {
      await server.checkForErrors(this.input).then(res => {
        this.output = res.data.message;
        this.lexical_errors = res.data.lexical_errors;
        this.syntax_errors = res.data.syntax_errors;
      });
    }
  }
}

</script>

<style scoped>
textarea {
  font-size: 10px;
}
#syntax-errors {
  color: red;
}

button {
  width: 100%;
}
</style>
