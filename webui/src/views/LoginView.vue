
<script> 
    import ErrorMsg from '../components/ErrorMsg.vue'
    
    export default{
      props: ['msg'],
        components: {
          ErrorMsg
        } ,
        data: function(){

            return {
                errMsg: null,
                username: "",
                userId: ""
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
                    let userId = response.data.identifier
                    // console.log(response.data.identifier)
                    // TODO go to homepage (or profile here)

                    // alert with the identifier
                    alert(`Welcome to WasaPHOTO, this is your id: ${userId}`)

                    this.$axios.defaults.headers.common['Authorization'] = `${response.data.identifier}`

                } catch(e) {
                    if (e.response.status === 400) {
                        this.errMsg = `Invalid username ${e.response.data}`
                        alert(this.errMsg)
                    }
                  }
            }
            },
            closeErrorMsg(){
                this.errMsg = null
            },
            mounted(){
                // TODO check what does mounted do 
            }
    }
</script>

<template>
  <div class="d-flex justify-content-center flex-column align-items-center">
    <!-- <ErrorMsg v-if="errMsg" :msg="errMsg"></ErrorMsg> -->
    <div class="d-flex justify-content-center flex-column align-items-center" style="width: 40%"> 
   	  <svg class="feather mt-0" style="width: 15%; height: 15%;"><use href="/feather-sprite-v4.29.0.svg#camera"/></svg>
    </div>
    <!-- <img class="mt-4 align-items-center" src="../assets/icons/w-button-icon.svg" alt="" width="80" height="80"> -->
    <p class="h4 mb-2 fw-bold">LOGIN</p>
    <div class="container" style ="width: 50%;">
      <form class="d-flex flex-column form-signin">
        <input type="text" id="username" v-model="username" class="form-control" placeholder="username" @input="closeErrorMsg">
        <button class="btn btn-primary mt-2" @click="doLogin">Let's go</button>
      </form>
    </div>
  </div>

  
</template>

<style></style>