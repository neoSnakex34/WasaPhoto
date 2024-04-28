
<script>
    export default{
        // TODO comment id will not be visualized, just used to remove 
        props: ["commentingUserId", "commentId",  "username", "body", "date", "photoOwnerId", "loggedUserId"], // add others 
        created () {
            console.log(this.commentingUserId + " "+ this.body) 
            
            // console.log(this.loggedUserId)
        }
    }
    
</script>

<!-- fix styling; add functions-->
<template>
    <div class="container pt-2 border-bottom">
        <div class="d-flex align-items-center"> <!-- add others -->
            <!-- commenting username should be and not logged in username, warning -->
            <h5 class="fw-bold"><strong>{{ this.username }}</strong></h5> <!--- username of commentor-->
            <!-- <h6 class="fw-bold text-muted"><strong>{{ this.commentingId }}</strong></h6> -->
            <span class="ms-4">{{ this.date }}</span>
            <button @click="$emit('delete-comment-event', commentId)" v-if="this.loggedUserId === this.photoOwnerId || this.loggedUserId === this.commentingUserId" class="btn btn-danger rounded-pill fw-bold ms-auto">x</button>
        </div>
        <div>
            <span>{{ this.body }}</span>
        </div>

    </div>
</template>