<script>
export default {
    // comment id will not be visualized, just used to remove 
    props: ["commentingUserId", "commentId", "username", "body", "date", "photoOwnerId", "loggedUserId"], // add others 
    created() {
        console.log(this.commentingUserId + " " + this.body)
    }
}

</script>

<!-- fix styling; add functions-->
<template>
    <div class="container pt-2 border-bottom">
        <div class="d-flex align-items-center"> <!-- add others -->
            <!-- commenting username should be and not logged in username, warning -->
            <div class="d-flex flex-column align-items-left">
                <div style="font-size: 14px;" class="pt-1">{{ this.date }}</div>
                <div style="font-size: 18px;" ><strong>{{ this.username }}</strong></div> <!--- username of commentor-->

            </div> 


            <button @click="$emit('delete-comment-event', commentId)"
                v-if="this.loggedUserId === this.photoOwnerId || this.loggedUserId === this.commentingUserId"
                class="delete-btn fw-bold ms-auto">x</button>
        </div>
        
        <div class="d-flex body-section pt-2">{{ this.body }}</div>

    </div>
</template>

<style>
  
    .body-section {
        font-size: 15px;
        font-style: italic;
    }

    .delete-btn {
    border: none; 
    outline: none; 
    color: white;
    background-color: transparent;  
    font-size: 25px;
    -webkit-text-stroke : 1px black;
    }
</style>