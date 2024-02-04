<script>
  export default {
    name: 'Modal',
    methods: {
      close() {
        this.$emit('close');
      },
      
    },
    mounted() {
      console.log("Modal created, id = " + this.id);
    },
  };
</script>
<template>
  <div class="modal fade" title="ModalWindow" tabindex="-1" role="dialog" aria-labelledby="Modal" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered modal-lg" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="ModalTitle">
            <!---->
              <slot name="header">
                This is the default title!
              </slot>
          </h5>
          <button type="button" class="close" data-bs-dismiss="modal" aria-label="Close">
            <span aria-hidden="true">&times;</span>
          </button>
        </div>
        <div class="modal-body justify-content-center align-items-center" style="overflow-x: auto; overflow-y: auto;">
          <h5 id="ModalBody">
            <!---->
            <slot name="body">
              Nothing in modal body
            </slot> 
          </h5>
              
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="close" data-bs-dismiss="modal">Close</button>
          <button type="button" class="btn btn-primary">Save changes</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
.modal-fade {
  display: table-cell;
  vertical-align: middle;
}

.modal-content {
  box-sizing: border-box;
  width: 6000px;
  margin: 0px auto;
  padding: 20px 30px;
  background-color: #fff;
  border-radius: 2px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.33);
  transition: all 0.3s ease;
}

.modal-header h3 {
  margin-top: 0;
  color: #42b983;
}

.modal-body {
  margin: 20px 0;
}

.modal-default-button {
  float: right;
}

/*
 * The following styles are auto-applied to elements with
 * transition="modal" when their visibility is toggled
 * by Vue.js.
 *
 * You can easily play with the modal transition by editing
 * these styles.
 */

.modal-enter-from {
  opacity: 0;
}

.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  -webkit-transform: scale(1.1);
  transform: scale(1.1);
}
</style>