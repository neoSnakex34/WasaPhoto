<script>
export default {
    props: ['msg'],
    data: function () {
        return {
            errMsg: null,
            username: localStorage.getItem('username'),
            userId: localStorage.getItem('userId'),
            clicked: false,
            newUsername: ''
        }
    },
    methods: {
        async setMyUsername(newUsername) {
            try {
                let response = await this.$axios.put(`/users/${this.userId}/username`, 
                    JSON.stringify(newUsername),
                    {
                    headers: {
                        Authorization: this.userId,
                        'Content-Type': 'application/json'
                    }    
                })

                localStorage.setItem('username', newUsername)
            }catch(e) {

            }
        }, 
        refresh() {
            this.username = localStorage.getItem('username')
            // TODO handle a real refresh
        },

        toggleEditing(newUsername) {
            if (this.clicked) {
                alert(this.clicked)
                this.setMyUsername(newUsername)
                this.refresh()
            }
            this.clicked = !this.clicked
        }
    },

}
</script>

<template>
    <div class="d-flex justify-content-left flex-wrap align-items-baseline pt-3 pb-2 mb-3 border-bottom">
        <h1 class="h1"><strong>{{ this.username }}</strong>'s profile</h1>
        <h6 class="text-muted ms-2">(<strong>{{ this.userId }}</strong>)</h6>
        <span class="d-flex justify-content-between align-items-center ms-auto">  
            <input v-if="clicked" v-model="newUsername" type="text" class="form-control form-control-lg me-4" placeholder="new username"/>
            <button class="btn btn-primary btn-lg rounded-pill  fw-bold" @click="toggleEditing(newUsername)">{{ clicked ? 'Confirm' : 'Edit username' }}</button>
        </span>
       


    </div>
</template>
