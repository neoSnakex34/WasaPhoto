
<script> 
import Comment from '../components/Comment.vue'

export default {
    components: {
        Comment
    },
    props: ['src', 'uploader', 'date', 'likes','liked'], // some ID wont be visualized
    data: function(){
        return {
            liked: false,
            likes: 0,
           
        }
    },
    mounted(){
        alert(this.src)
    },
    
}
</script>


<!--- href for visualizing photo in another windows? what about -->
<template>
    
        <div class="card flex-grow-1">
            <img :src="src" class="card-img-top" />
            <div class="card-body">
                <p class="card-text">Uploader: {{ this.uploader }}</p> <!-- username href to profile -->
                <p class="card-text">Date: {{ this.date }}</p>
                <p class="card-text">Likes: {{ this.likes }}</p>
                <!-- this should become an icon probably -->
                <p class="card-text">Liked: {{ this.liked }}</p>
            </div>
            
        </div>
        <div class="input-group rounded pt-1">
                <input class="form-control form-control-lg" type="text" placeholder="Comment"/>
                <button class="btn btn-success btn-lg fw-bold" type="button">Comment</button>
        </div>

          <!-- change accordingly with photo max dimension, must be set-->
        <div class="overflow-auto  pt-2 pb-5 mb-5" style="max-height: 200px;">
            <Comment 
                commentingId="1"
                username="rei"
                body="This is a comment"
                date="2021-10-10"
                />
             <Comment 
                commentingId="1"
                username="rei"
                body="This is a comment"
                date="2021-10-10"
                /> <Comment 
                commentingId="1"
                username="rei"
                body="This is a comment"
                date="2021-10-10"
                /> <Comment 
                commentingId="1"
                username="rei"
                body="This is a comment"
                date="2021-10-10"
                /> <Comment 
                commentingId="1"
                username="rei"
                body="This is a comment"
                date="2021-10-10"
                /> <Comment 
                commentingId="1"
                username="rei"
                body="This is a comment"
                date="2021-10-10"
                /> <Comment 
                commentingId="1"
                username="rei"
                body="This is a comment"
                date="2021-10-10"
                /> <Comment 
                commentingId="1"
                username="rei"
                body="This is a comment"
                date="2021-10-10"
                /> <Comment 
                commentingId="1"
                username="rei"
                body="This is a comment"
                date="2021-10-10"
                />

        </div>
</template>
