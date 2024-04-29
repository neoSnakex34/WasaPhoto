
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
                    localStorage.setItem('userId', userId)
                    localStorage.setItem('username', this.username)
                    localStorage.setItem('authorization', userId)
                   
                    // [ ] go to homepage (or profile here)
                    console.log(response)
                    this.$router.push('/home')
                    


                    this.$axios.defaults.headers.common['Authorization'] = localStorage.getItem('authorization')
                    // alert with the identifier
                    // alert(`Welcome to WasaPHOTO, this is your id: ${userId}`)

                    

                } catch(e) {
                  // TODO check other errors
                    if (e.response.data) {
                        this.errMsg = `${e.response.data}`
                        // alert(this.errMsg)
                    } else{alert (e)}
                  }
            }
            },
        
          
            mounted(){
                // TODO check what does mounted do 
            }
    }
</script>

<template>
  <div class="d-flex justify-content-center flex-column align-items-center">
    
    
    <div class="d-flex justify-content-center flex-column align-items-center" style="width: 40%"> 
   	  <svg class="feather mt-4" style="min-width: 100px; min-height: 100px; width: 12%; height: 12%;"><use href="/feather-sprite-v4.29.0.svg#camera"/></svg>
    </div>
    <!-- <img class="mt-4 align-items-center" src="../assets/icons/w-button-icon.svg" alt="" width="80" height="80"> -->
    <p class="h1 mb-2 fw-bold">WasaPHOTO</p>
    <p class="h3 mb-2 fw-bold">LOGIN</p>
    <div class="container mt-4" style ="min-width: 400px; min-height: 50px; width: 50%; height: 60%;">
      <div class="d-flex flex-column">
        <input type="text" title="only lowercase alphanumeric, min3 max12" v-model="username" class="form-control form-control-lg" placeholder="username" style="outline: 2px solid lightblue;">
        <button class="btn btn-primary mt-2 fw-bold btn-lg" @click="doLogin">Let's go</button>
      </div>
    </div>
  </div>

  <div class ="d-flex justify-content-center align-items-center pt-3 fw-bold"> 
    <!-- : it's the shorthand for v-bind -->
    <ErrorMsg v-if="errMsg" :msg="errMsg" />
    <button v-if="errMsg" class="flex-column btn btn-danger ms-3 mb-3 fw-bold" @click="errMsg = null">Got it</button>
  </div>
  
</template>

<style></style>