document.addEventListener('DOMContentLoaded', function() {
    // App state
    const state = {
        currentQuiz: null,
        currentQuestion: null,
        isEditingQuestion: false,
        optionCounter: 2, // Start with 2 options by default
    };

    // API Base URL
    const API_BASE_URL = '/api';

    // Elements
    const tabButtons = document.querySelectorAll('.tab-btn');
    const tabPanes = document.querySelectorAll('.tab-pane');
    const quizzesList = document.getElementById('quizzes-list');
    const quizForm = document.getElementById('quiz-form');
    const editQuizForm = document.getElementById('edit-quiz-form');
    const backBtn = document.querySelector('.back-btn');
    const backToDetailBtns = document.querySelectorAll('.back-to-detail-btn');
    const editQuizBtn = document.getElementById('edit-quiz-btn');
    const deleteQuizBtn = document.getElementById('delete-quiz-btn');
    const addQuestionBtn = document.getElementById('add-question-btn');
    const addQuestionForm = document.getElementById('add-question-form');
    const addOptionBtn = document.getElementById('add-option-btn');
    const optionsList = document.getElementById('options-list');
    const questionSubmitText = document.getElementById('question-submit-text');
    const questionFormTitle = document.getElementById('question-form-title');
    const modal = document.getElementById('modal');
    const modalTitle = document.getElementById('modal-title');
    const modalBody = document.getElementById('modal-body');
    const modalConfirm = document.getElementById('modal-confirm');
    const modalCancel = document.getElementById('modal-cancel');
    const closeModal = document.querySelector('.close-modal');

    // Event Listeners
    tabButtons.forEach(btn => {
        btn.addEventListener('click', () => {
            const tabId = btn.dataset.tab;
            switchTab(tabId);
        });
    });

    quizForm.addEventListener('submit', createQuiz);
    editQuizForm.addEventListener('submit', updateQuiz);
    backBtn.addEventListener('click', goBackToQuizzes);
    
    backToDetailBtns.forEach(btn => {
        btn.addEventListener('click', goBackToQuizDetail);
    });
    
    editQuizBtn.addEventListener('click', showEditQuizForm);
    deleteQuizBtn.addEventListener('click', confirmDeleteQuiz);
    addQuestionBtn.addEventListener('click', showAddQuestionForm);
    addQuestionForm.addEventListener('submit', handleQuestionSubmit);
    addOptionBtn.addEventListener('click', addNewOption);
    
    closeModal.addEventListener('click', hideModal);
    modalCancel.addEventListener('click', hideModal);

    // Initialize app
    loadQuizzes();

    // Switch tab function
    function switchTab(tabId) {
        tabButtons.forEach(btn => {
            btn.classList.remove('active');
            if (btn.dataset.tab === tabId) {
                btn.classList.add('active');
            }
        });

        tabPanes.forEach(pane => {
            pane.classList.remove('active');
            if (pane.id === tabId) {
                pane.classList.add('active');
            }
        });
    }

    // Load quizzes from API
    async function loadQuizzes() {
        try {
            const response = await fetch(`${API_BASE_URL}/quizzes`);
            if (!response.ok) throw new Error('Failed to fetch quizzes');
            
            const quizzes = await response.json();
            renderQuizzes(quizzes);
        } catch (error) {
            console.error('Error loading quizzes:', error);
            quizzesList.innerHTML = `
                <div class="quiz-card">
                    <h3><i class="fas fa-exclamation-triangle"></i> Error</h3>
                    <p>Failed to load quizzes. Please try again later.</p>
                </div>
            `;
        }
    }

    // Render quizzes to the DOM
    function renderQuizzes(quizzes) {
        if (quizzes.length === 0) {
            quizzesList.innerHTML = `
                <div class="quiz-card">
                    <h3>No Quizzes Found</h3>
                    <p>Create your first quiz by clicking the "Create Quiz" tab.</p>
                </div>
            `;
            return;
        }

        quizzesList.innerHTML = '';
        quizzes.forEach(quiz => {
            const card = document.createElement('div');
            card.className = 'quiz-card';
            card.innerHTML = `
                <h3>${quiz.title}</h3>
                <p>${quiz.questions ? quiz.questions.length : 0} questions</p>
                <div class="quiz-card-footer">
                    <span>Created: ${formatDate(quiz.created_at)}</span>
                    <button class="btn primary small view-quiz" data-id="${quiz.uuid}">
                        <i class="fas fa-eye"></i> View
                    </button>
                </div>
            `;
            quizzesList.appendChild(card);

            // Add event listener to view button
            card.querySelector('.view-quiz').addEventListener('click', () => {
                loadQuizDetail(quiz.uuid);
            });
        });
    }

    // Create a new quiz
    async function createQuiz(e) {
        e.preventDefault();
        const title = document.getElementById('quiz-title').value.trim();
        
        if (!title) {
            alert('Please enter a quiz title');
            return;
        }

        try {
            const response = await fetch(`${API_BASE_URL}/quizzes`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ title })
            });

            if (!response.ok) throw new Error('Failed to create quiz');
            
            const quiz = await response.json();
            quizForm.reset();
            loadQuizDetail(quiz.uuid);
        } catch (error) {
            console.error('Error creating quiz:', error);
            alert('Failed to create quiz. Please try again.');
        }
    }

    // Load quiz detail
    async function loadQuizDetail(quizId) {
        try {
            const response = await fetch(`${API_BASE_URL}/quizzes/${quizId}`);
            if (!response.ok) throw new Error('Failed to fetch quiz details');
            
            const quiz = await response.json();
            state.currentQuiz = quiz;
            
            renderQuizDetail(quiz);
            switchTab('quiz-detail');
        } catch (error) {
            console.error('Error loading quiz details:', error);
            alert('Failed to load quiz details. Please try again.');
        }
    }

    // Render quiz detail
    function renderQuizDetail(quiz) {
        document.getElementById('quiz-detail-title').textContent = quiz.title;
        
        // Load questions
        loadQuizQuestions(quiz.uuid);
    }

    // Load quiz questions
    async function loadQuizQuestions(quizId) {
        const questionsList = document.getElementById('questions-list');
        questionsList.innerHTML = '<div class="loading-spinner"></div>';
        
        try {
            const response = await fetch(`${API_BASE_URL}/quizzes/${quizId}/questions`);
            if (!response.ok) throw new Error('Failed to fetch questions');
            
            const questions = await response.json();
            renderQuestions(questions);
        } catch (error) {
            console.error('Error loading questions:', error);
            questionsList.innerHTML = `
                <div class="question-card">
                    <p><i class="fas fa-exclamation-triangle"></i> Failed to load questions. Please try again.</p>
                </div>
            `;
        }
    }

    // Render questions
    function renderQuestions(questions) {
        const questionsList = document.getElementById('questions-list');
        
        if (questions.length === 0) {
            questionsList.innerHTML = `
                <div class="question-card">
                    <p>No questions yet. Add your first question!</p>
                </div>
            `;
            return;
        }

        questionsList.innerHTML = '';
        questions.forEach((question, index) => {
            const card = document.createElement('div');
            card.className = 'question-card';
            
            // Generate options HTML
            let optionsHTML = '';
            if (question.options && question.options.length > 0) {
                optionsHTML = '<div class="options-list">';
                question.options.forEach(option => {
                    const isCorrect = option.is_correct ? 'correct' : '';
                    optionsHTML += `
                        <div class="option-card ${isCorrect}">
                            ${option.text}
                            ${option.is_correct ? '<span class="badge">✓ Correct</span>' : ''}
                        </div>
                    `;
                });
                optionsHTML += '</div>';
            }
            
            card.innerHTML = `
                <div class="question-header">
                    <div class="question-content">
                        <h3>Q${index + 1}: ${question.text}</h3>
                        <div class="question-meta">
                            <span><i class="fas fa-clock"></i> ${question.time_limit}s</span>
                            <span><i class="fas fa-star"></i> ${question.points} points</span>
                        </div>
                    </div>
                    <div class="question-actions">
                        <button class="btn secondary small edit-question" data-id="${question.uuid}">
                            <i class="fas fa-edit"></i> Edit
                        </button>
                        <button class="btn danger small delete-question" data-id="${question.uuid}">
                            <i class="fas fa-trash"></i> Delete
                        </button>
                    </div>
                </div>
                ${optionsHTML}
            `;
            
            questionsList.appendChild(card);
            
            // Add event listeners
            card.querySelector('.edit-question').addEventListener('click', () => {
                editQuestion(question);
            });
            
            card.querySelector('.delete-question').addEventListener('click', () => {
                confirmDeleteQuestion(question.uuid);
            });
        });
    }

    // Show edit quiz form
    function showEditQuizForm() {
        document.getElementById('edit-quiz-title').value = state.currentQuiz.title;
        switchTab('edit-quiz');
    }

    // Update quiz
    async function updateQuiz(e) {
        e.preventDefault();
        const title = document.getElementById('edit-quiz-title').value.trim();
        
        if (!title) {
            alert('Please enter a quiz title');
            return;
        }

        try {
            const response = await fetch(`${API_BASE_URL}/quizzes/${state.currentQuiz.uuid}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ title })
            });

            if (!response.ok) throw new Error('Failed to update quiz');
            
            const quiz = await response.json();
            state.currentQuiz = quiz;
            
            renderQuizDetail(quiz);
            switchTab('quiz-detail');
        } catch (error) {
            console.error('Error updating quiz:', error);
            alert('Failed to update quiz. Please try again.');
        }
    }

    // Show add question form
    function showAddQuestionForm() {
        state.isEditingQuestion = false;
        state.currentQuestion = null;
        
        // Reset form
        addQuestionForm.reset();
        
        // Reset options (keep only two default options)
        optionsList.innerHTML = `
            <div class="option-item">
                <input type="text" name="option_text_1" placeholder="Option text" required>
                <label class="checkbox-container">
                    <input type="radio" name="correct_option" value="0" checked>
                    <span class="checkmark"></span>
                    Correct
                </label>
                <button type="button" class="btn danger small remove-option"><i class="fas fa-times"></i></button>
            </div>
            <div class="option-item">
                <input type="text" name="option_text_2" placeholder="Option text" required>
                <label class="checkbox-container">
                    <input type="radio" name="correct_option" value="1">
                    <span class="checkmark"></span>
                    Correct
                </label>
                <button type="button" class="btn danger small remove-option"><i class="fas fa-times"></i></button>
            </div>
        `;
        
        // Add event listeners to remove buttons
        addRemoveOptionListeners();
        
        state.optionCounter = 2;
        questionFormTitle.textContent = 'Add Question';
        questionSubmitText.textContent = 'Add Question';
        
        switchTab('question-form');
    }

    // Edit question
    function editQuestion(question) {
        state.isEditingQuestion = true;
        state.currentQuestion = question;
        
        // Fill form with question data
        document.getElementById('question-text').value = question.text;
        document.getElementById('question-time-limit').value = question.time_limit;
        document.getElementById('question-points').value = question.points;
        
        // Create options
        optionsList.innerHTML = '';
        state.optionCounter = 0;
        
        if (question.options && question.options.length > 0) {
            question.options.forEach((option, index) => {
                addNewOption(null, option.text, option.is_correct);
            });
        } else {
            // Add two empty options if none exist
            addNewOption();
            addNewOption();
        }
        
        questionFormTitle.textContent = 'Edit Question';
        questionSubmitText.textContent = 'Save Changes';
        
        switchTab('question-form');
    }

    // Handle question form submit (add or edit)
    async function handleQuestionSubmit(e) {
        e.preventDefault();
        
        const questionText = document.getElementById('question-text').value.trim();
        const timeLimit = parseInt(document.getElementById('question-time-limit').value);
        const points = parseInt(document.getElementById('question-points').value);
        
        if (!questionText) {
            alert('Please enter question text');
            return;
        }
        
        // Get options
        const options = [];
        const optionItems = optionsList.querySelectorAll('.option-item');
        const selectedCorrectOption = document.querySelector('input[name="correct_option"]:checked').value;
        
        optionItems.forEach((item, index) => {
            const optionText = item.querySelector('input[type="text"]').value.trim();
            if (optionText) {
                options.push({
                    text: optionText,
                    is_correct: index.toString() === selectedCorrectOption
                });
            }
        });
        
        if (options.length < 2) {
            alert('Please add at least two options');
            return;
        }
        
        const questionData = {
            text: questionText,
            time_limit: timeLimit,
            points: points,
            options: options
        };
        
        try {
            let response;
            
            if (state.isEditingQuestion) {
                // Update existing question
                response = await fetch(`${API_BASE_URL}/questions/${state.currentQuestion.uuid}`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(questionData)
                });
            } else {
                // Add new question
                response = await fetch(`${API_BASE_URL}/quizzes/${state.currentQuiz.uuid}/questions`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(questionData)
                });
            }

            if (!response.ok) throw new Error('Failed to save question');
            
            // Reset form and go back to quiz detail
            addQuestionForm.reset();
            loadQuizQuestions(state.currentQuiz.uuid);
            switchTab('quiz-detail');
        } catch (error) {
            console.error('Error saving question:', error);
            alert('Failed to save question. Please try again.');
        }
    }

    // Add new option to the question form
    function addNewOption(e, text = '', isCorrect = false) {
        if (e) e.preventDefault();
        
        state.optionCounter++;
        const newOption = document.createElement('div');
        newOption.className = 'option-item';
        
        newOption.innerHTML = `
            <input type="text" name="option_text_${state.optionCounter}" placeholder="Option text" required value="${text}">
            <label class="checkbox-container">
                <input type="radio" name="correct_option" value="${state.optionCounter - 1}" ${isCorrect ? 'checked' : ''}>
                <span class="checkmark"></span>
                Correct
            </label>
            <button type="button" class="btn danger small remove-option"><i class="fas fa-times"></i></button>
        `;
        
        optionsList.appendChild(newOption);
        
        // Add event listener to remove button
        newOption.querySelector('.remove-option').addEventListener('click', removeOption);
    }

    // Add event listeners to remove option buttons
    function addRemoveOptionListeners() {
        document.querySelectorAll('.remove-option').forEach(btn => {
            btn.addEventListener('click', removeOption);
        });
    }

    // Remove an option
    function removeOption(e) {
        e.preventDefault();
        
        const optionItems = optionsList.querySelectorAll('.option-item');
        if (optionItems.length <= 2) {
            alert('You need at least two options');
            return;
        }
        
        e.target.closest('.option-item').remove();
        
        // Update radio button values to be consecutive
        const radioButtons = optionsList.querySelectorAll('input[type="radio"]');
        radioButtons.forEach((radio, index) => {
            radio.value = index;
        });
    }

    // Confirm delete quiz
    function confirmDeleteQuiz() {
        modalTitle.textContent = 'Delete Quiz';
        modalBody.innerHTML = `
            <p>Are you sure you want to delete the quiz "${state.currentQuiz.title}"?</p>
            <p>This action cannot be undone.</p>
        `;
        modalConfirm.onclick = deleteQuiz;
        showModal();
    }

    // Delete quiz
    async function deleteQuiz() {
        hideModal();
        
        try {
            const response = await fetch(`${API_BASE_URL}/quizzes/${state.currentQuiz.uuid}`, {
                method: 'DELETE'
            });

            if (!response.ok) throw new Error('Failed to delete quiz');
            
            goBackToQuizzes();
            loadQuizzes();
        } catch (error) {
            console.error('Error deleting quiz:', error);
            alert('Failed to delete quiz. Please try again.');
        }
    }

    // Confirm delete question
    function confirmDeleteQuestion(questionId) {
        modalTitle.textContent = 'Delete Question';
        modalBody.innerHTML = `
            <p>Are you sure you want to delete this question?</p>
            <p>This action cannot be undone.</p>
        `;
        modalConfirm.onclick = () => deleteQuestion(questionId);
        showModal();
    }

    // Delete question
    async function deleteQuestion(questionId) {
        hideModal();
        
        try {
            const response = await fetch(`${API_BASE_URL}/questions/${questionId}`, {
                method: 'DELETE'
            });

            if (!response.ok) throw new Error('Failed to delete question');
            
            loadQuizQuestions(state.currentQuiz.uuid);
        } catch (error) {
            console.error('Error deleting question:', error);
            alert('Failed to delete question. Please try again.');
        }
    }

    // Show modal
    function showModal() {
        modal.style.display = 'block';
    }

    // Hide modal
    function hideModal() {
        modal.style.display = 'none';
    }

    // Go back to quizzes list
    function goBackToQuizzes() {
        state.currentQuiz = null;
        switchTab('quizzes');
    }

    // Go back to quiz detail
    function goBackToQuizDetail() {
        switchTab('quiz-detail');
    }

    // Format date
    function formatDate(dateString) {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric'
        });
    }

    // Initialize event listeners for dynamically added elements
    addRemoveOptionListeners();

    // Close modal when clicking outside
    window.onclick = function(event) {
        if (event.target == modal) {
            hideModal();
        }
    }
});