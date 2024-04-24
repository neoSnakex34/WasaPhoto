<script>
export default {
    data: function () {
        return {
            profile: {
                username: "",
                userId: "",
                following: 0,
                followedCounter: 0,
                photoCounter: 0,
                photos: [],

            },

            servedPhotos: [],
            banned: false,

        }
    }
}
</script>

<template>
    <div class="pt-3 pb-3 border-bottom">

        <div class="d-flex">
            <!-- username and id -->
            <div class="d-flex align-items-baseline ps-1">
                <h3 class="h3"><strong>username</strong></h3>
                <h6 class="text-muted ms-2">(<strong>userid</strong>)</h6>
            </div>

            <div class="d-flex ms-auto pe-1">
                <!-- counters -->
                <div class="text-center me-4">
                    <div class="fw-bold">Following</div>
                    <div class="fw-bold">0</div>
                </div>
                <div class="text-center ps-4 pe-4 me-4">
                    <div class="fw-bold">Followed</div>
                    <div class="fw-bold">0</div>
                </div>
                <div class="text-center">
                    <div class="fw-bold">Photos</div>
                    <div class="fw-bold">0</div>
                </div>
            </div>

        </div>



    </div>
    <div class="container pt-3 pb-3 d-flex align-items-center justify-content-center" style="width: 90%;">


    </div>


    <div class="border-bottom"></div>

    <!-- photos v if not banned -->
    <div class="container pt-4 pb-4" style="width: 60%;">
        <div>

            <!--- in stream uploader  will need to be passed from struct -->


        </div>
    </div>
</template>
