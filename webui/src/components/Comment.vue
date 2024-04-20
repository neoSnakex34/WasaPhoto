
<script>
    export default{
        // TODO comment id will not be visualized, just used to remove 
        props: ['commentId', 'body','username', 'commentingId', 'date']
    }
</script>

<!-- fix styling; add functions-->
<template>
    <div class="container pt-2 border-bottom">
        <div class="d-flex"> <!-- add others -->
            <h5 class="fw-bold"><strong>{{ this.username }}</strong></h5> <!--- username of commentor-->
            <h6 class="fw-bold text-muted"><strong>{{ this.commentingId }}</strong></h6>
            <span>{{ this.date }}</span>
        </div>
        <div>
            <span>{{ this.body }}</span>
        </div>
    </div>
</template>