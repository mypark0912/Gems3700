<!-- components/LoadingModal.vue -->
<template>
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="isOpen" class="modal-overlay">
          <div class="modal-box">
            <div class="modal-content">
              <div class="spinner"></div>
              <h3>{{ title }}</h3>
              <p>{{ message }}</p>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </template>
  
  <script setup>
  defineProps({
    isOpen: Boolean,
    title: String,
    message: String
  });
  </script>
  
  <style scoped>
  .modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.8);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
  }
  
  .modal-box {
    background: #1e1e1e;
    border-radius: 12px;
    padding: 40px;
    min-width: 320px;
  }
  
  .modal-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;
  }
  
  .spinner {
    width: 50px;
    height: 50px;
    border: 4px solid #374151;
    border-top-color: #3b82f6;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
  
  h3 {
    color: white;
    font-size: 18px;
    font-weight: 600;
    margin: 0;
  }
  
  p {
    color: #9ca3af;
    font-size: 14px;
    margin: 0;
    text-align: center;
  }
  </style>