<script>
import Comment from '../components/Comment.vue'

export default {
    data: function () {
        return {
            likerId: localStorage.getItem('userId')
        }
    },
    components: {
        Comment
    },
    props: ['src', 'uploader', 'uploaderId', 'date', 'likes', 'liked', 'photoId'], // some ID wont be visualized
    methods:
    {
        async likePhoto() {
            try {
                let response = await this.$axios.put(`/users/${this.uploaderId}/photos/${this.photoId}/likes/${this.likerId}`,
                    {
                        headers: {
                            // Authorization: localStorage.getItem('authorization'),
                            'Content-Type': 'application/json'
                        }
                    }

                )
                if (response.status === 200) {
                    this.$emit('like');
                }
                // this.likes++
            } catch (e) {
                if (e.response.data) {
                    alert(e.response.data)
                } else {
                    alert(e)
                }
            }
        },

        async unlikePhoto() {

        }

    },
    mounted() {
        // alert(this.src)
    },

}
</script>


<!--- href for visualizing photo in another windows? what about -->
<template>

    <div class="card flex-grow-1">
        <img :src="src" class="card-img-top" />
        <div class="card-body d-flex flex-column">
        
            <div class="d-flex align-items-center justify-content-center justify-content-between pt-2 pb-2"
                style="width: 60%; margin: auto;">
                <div class="d-flex flex-column justify-content-between">
                    <p class="card-text"><strong>UPLOADER</strong>: {{ this.uploader }}</p>
                    <!-- username href to<strong> profile -->
                    <p class="card-text"><strong>DATE</strong> {{ this.date }}</p>
                    <p class="card-text"><strong>LIKES</strong>: {{ this.likes }}</p>
                    <!-- this should beco<strong>me an icon probably -->
                    <p class="card-text"><strong>LIKED</strong>: {{ this.liked }}</p>

                </div>
                <button class="ms-4 d-flex btn rounded-pill btn-primary fw-bold" @click="likePhoto()">Like</button>
                <button class="d-flex btn rounded-pill btn-danger fw-bold" @click="unlikePhoto()">Unlike</button>

            </div>
        </div>

    </div>
    <div class="input-group rounded pt-1">
        <input class="form-control form-control-lg" type="text" placeholder="Comment" />
        <button class="btn btn-success btn-lg fw-bold" type="button">Comment</button>
    </div>

    <!-- change accordingly with photo max dimension, must be set-->
    <div class="overflow-auto  pt-2 pb-5 mb-5" style="max-height: 200px;">
        <Comment commentingId="1" username="rei" body="This is a comment" date="2021-10-10" />
        <Comment commentingId="1" username="rei" body="This is a comment" date="2021-10-10" />
        <Comment commentingId="1" username="rei" body="This is a comment" date="2021-10-10" />
        <Comment commentingId="1" username="rei" body="This is a comment" date="2021-10-10" />
        <Comment commentingId="1" username="rei" body="This is a comment" date="2021-10-10" />
        <Comment commentingId="1" username="rei" body="This is a comment" date="2021-10-10" />
        <Comment commentingId="1" username="rei" body="This is a comment" date="2021-10-10" />
        <Comment commentingId="1" username="rei" body="This is a comment" date="2021-10-10" />
        <Comment commentingId="1" username="rei" body="This is a comment" date="2021-10-10" />

    </div>
</template>
