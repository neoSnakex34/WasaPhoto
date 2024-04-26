<script>
import Comment from '../components/Comment.vue'

export default {
    data: function () {
        return {
            likerId: localStorage.getItem('userId'),
            commentBody: '',
            comments: [
                /* 
                    commentId 
                    userId // id of commentor
                    username // username of commentor
                    photoId 
                    body
                    date
                */
            ]
        }
    },
    components: {
        Comment
    },
    props: ['src', 'uploader', 'uploaderId', 'date', 'likes', 'liked', 'photoId', 'delete', 'guest'], // some ID wont be visualized
    methods: {

        

        toggleDelete() {
            this.delete = !this.delete
        },

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


        async unlikePhoto(){
            try{
                let response = await this.$axios.delete(`/users/${this.uploaderId}/photos/${this.photoId}/likes/${this.likerId}`,
                {
                    headers: {
                        
                        'Content-Type': 'application/json'
                    }
                }
                )

                if(response.status === 200){
                    this.$emit('unlike');
                }
                
            }catch(e){
                if (e.response.data) {
                    alert(e.response.data)
                } else {
                    alert(e)
                }
            }
        },

        async commentPhoto() {

            try {
                alert(" called ")
                let response = await this.$axios.post(`/users/${this.uploaderId}/photos/${this.photoId}/comments`,
                    {
                        body: this.commentBody
                    },
                    {

                        headers: {
                            requestor: localStorage.getItem('userId'),
                            'Content-Type': 'application/json'
                        }
                    })

            } catch (e) {
                if (e.response.data) {
                    alert(e.response.data)
                } else {
                    alert(e)
                }

            }

            // empyting the comment field
            this.commentBody = ''
        },

    },

  

    
    mounted() {
        // alert(this.src)
    },

}
</script>


<!--- href for visualizing photo in another windows? what about -->
<template>

    <div class="d-flex flex-grow-1 pb-2">
        <button v-if="this.delete" @click="$emit('delete-event')" class="d-flex btn btn-danger flex-grow-1 justify-content-center fw-bold">DELETE</button>
    </div>
 
    <div class="card flex-grow-1 ">
        
       

        <img :src="src" class="card-img-top" />

            <!-- kebab button  v if profile personal (add a check) unused valued should be null in guestProfileView-->
            <button v-if="!guest" @click="$emit('toggle-delete')" class="position-absolute top-0 end-0 custom-button fw-bold">...</button>
            <!-- <button @click="toggleDelete" class="position-absolute top-0 end-0 custom-button fw-bold">...</button> -->
           
   
        <div class="card-body d-flex flex-column">
                  


            <div class="d-flex align-items-center justify-content-center justify-content-between pt-2 pb-2"
                style="width: 60%; margin: auto;">
                <div class="d-flex flex-column justify-content-between">
                    <p class="card-text"><strong>UPLOADER</strong>: {{ this.uploader }}</p>
                    <!-- username href to<strong> profile -->
                    <p class="card-text"><strong>DATE</strong> {{ this.date }}</p>
                    <p class="card-text"><strong>LIKES</strong>: {{ this.likes }}</p>
                    <!-- this should beco<strong>me an icon probably -->
                    <!-- <p class="card-text"><strong>LIKED</strong>: {{ this.liked }}</p> -->
                    
                    

                </div>
                
                
                <!-- heart icon using imported svgs from folder public  -->
                    <div  class="me-4" style="width: 10%;">
                        <div v-if="!this.liked">
                            <svg viewBox="0 0 512 512" xmlns="http://www.w3.org/2000/svg" @click="likePhoto()" style="fill: black; opacity: 75%;">
                                <path d="m458.4 64.3c-57.8-48.6-147.1-41.3-202.4 15-55.3-56.3-144.6-63.7-202.4-15-75.2 63.3-64.2 166.5-10.6 221.2l175.4 178.7c10 10.2 23.4 15.9 37.6 15.9 14.3 0 27.6-5.6 37.6-15.8l175.4-178.7c53.5-54.7 64.7-157.9-10.6-221.3zm-23.6 187.5-175.4 178.7c-2.4 2.4-4.4 2.4-6.8 0l-175.4-178.7c-36.5-37.2-43.9-107.6 7.3-150.7 38.9-32.7 98.9-27.8 136.5 10.5l35 35.7 35-35.7c37.8-38.5 97.8-43.2 136.5-10.6 51.1 43.1 43.5 113.9 7.3 150.8z"/>
                            </svg>
                        </div>
                        <div v-if="this.liked">
                            <svg viewBox="0 0 512 512" xmlns="http://www.w3.org/2000/svg" @click="unlikePhoto()" style="fill: blueviolet ">
                                <path d="m462.3 62.6c-54.8-46.7-136.3-38.3-186.6 13.6l-19.7 20.3-19.7-20.3c-50.2-51.9-131.8-60.3-186.6-13.6-62.8 53.6-66.1 149.8-9.9 207.9l193.5 199.8c12.5 12.9 32.8 12.9 45.3 0l193.5-199.8c56.3-58.1 53-154.3-9.8-207.9z"/>
                            </svg>
                        </div>
                    </div>

                <!--  deprecated, remove 
                <button class="ms-4 d-flex btn rounded-pill btn-primary fw-bold" @click="likePhoto()">Like</button>
                <button class="d-flex btn rounded-pill btn-danger fw-bold" @click="unlikePhoto()">Unlike</button> -->

            </div>
        </div>

    </div>
    <div class="input-group rounded pt-1">
        <!-- emit remove comment (if author == ecc ecc )-->
        <input  v-model="commentBody" class="form-control form-control-lg" type="text" placeholder="Comment" />
        <button @click="commentPhoto" class="btn btn-success btn-lg fw-bold" type="button">Comment</button>
    </div>

    <!-- change accordingly with photo max dimension, must be set-->
    <div class="overflow-auto  pt-2 pb-5 mb-5" style="max-height: 200px;">
        <Comment commentingId="1" username="rei" body="This is a comment" date="2021-10-10" />
        <Comment commentingId="1" username="rei" body="This is a comment" date="2021-10-10" />


    </div>

    
</template>

<style>
    .custom-button {
    border: none; 
    outline: none; 
    color: white;
   
    background-color: transparent;  
    font-size: 40px;
    text-shadow:  /* cross browser */
        -1px -1px 0 black,  /* shadow on every direction */
         1px -1px 0 black, /* simulates a stroke */
        -1px  1px 0 black, /* less aesthetic than webkit stroke but works */
         1px  1px 0 black;
    padding: 5px 30px; /* top bottom; left right */
}
</style>