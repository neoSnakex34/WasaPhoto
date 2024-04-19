<script>
    export default{
        props: ['msg'],
        data: function(){
            return {
                errMsg: null,
                username: localStorage.getItem('username'),
                userId: localStorage.getItem('userId')
            }
        },
        

    }
</script>

<template>
    <div class="d-flex justify-content-left flex-wrap align-items-baseline pt-3 pb-2 mb-3 border-bottom">
        <h1 class="h1"><strong>{{ this.username }}</strong>'s profile</h1>
        <h6 class="text-muted ms-2">(<strong>{{ this.userId }}</strong>)</h6>
        <!-- <div class="btn-toolbar mb-2 mb-md-0">
            <div class="me-4">
                <button type="button" class="btn btn-danger btn-lg fw-bold rounded-pill" @click="doLogout">
                    Logout
                </button>
            </div>
  
        </div> -->

    </div>
</template>
