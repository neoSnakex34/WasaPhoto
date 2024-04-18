
<script> 
    export default{
        components: {} ,
        data: function(){

            return {
                errMsg: null,
                username: "",
            }
        },

        methods: {
          async doLogin(){
                try {

                    let response = await this.$axios.post('/login', this.username, {
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    
                    })
                    console.log(response.data.identifier)
                    // TODO go to homepage (or profile here)

                    //
                    this.$axios.defaults.headers.common['Authorization'] = `${response.data.identifier}`

                } catch(e) {

                }
            }
            },
            mounted(){
                // TODO check what does mounted do 
            }
    }
</script>

<template>
  <div class ="input-group">
    <input type="text" id="username" v-model="username" class="form-contorl" placeholder="username">
      <div>
        <button class ="btn" type="button" @click="doLogin">Log me in!</button>
      </div> 
  </div>
</template>

<style scoped>
  .container{
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;

  }

  .input-group {
    display: flex;
    align-items: center;
  }

  .btn {
    margin-top: 10px;
    margin-left: 10px;
    align-items: center;
    background: black;
    color: blue;
  }
</style>