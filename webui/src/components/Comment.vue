
<script>
    export default{
        // TODO comment id will not be visualized, just used to remove 
        props: ["commentingId", "username", "body", "date"], // add others 
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
            <button class="btn btn-danger rounded-pill ms-auto">x</button>
        </div>
        <div>
            <span>{{ this.body }}</span>
        </div>

    </div>
</template>