function getPageFromUrl() {
    const url = new URL(window.location.href);
    return url.pathname.split('/').pop();
}

// 获取需要操作的DOM元素
const commentInput = $('#comment_input');
const parentIdInput = $('#parent-id_input');
const rootIdInput = $('#root-id_input');
const commentForm = $('#comment-form');
const commentList = $('#comment-list');
const submitButton = $('#submit-button');
const usernameInput = $('#username_input');
const websiteInput = $('#website_input');

// 当用户点击提交按钮时
submitButton.on('click', function() {
    // 获取用户的输入
    const comment = commentInput.val();
    const parentId = parentIdInput.val();
    const rootId = rootIdInput.val();
    const username = usernameInput.val();
    const website = websiteInput.val();
    if (comment.length===0){
        return
    }
    if (username.length===0){
        return
    }
    // 向后端API发送一个POST请求，包含用户的留言
    $.ajax({
        url: '/comments',
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
            page:getPageFromUrl(),
            username: username,
            website: website,
            content: comment,
            parentId: parseInt(parentId) || 0,
            rootId: parseInt(rootId)||0
        }),
        success: function(data) {
            // 清空输入框
            commentInput.val('');
            parentIdInput.val('');
            rootIdInput.val('');

            // 加载所有的留言
            loadComments();
        }
    });
});

function createCommentElement(comment) {
    // 创建一个新的div元素，用作留言容器
    const div = $('<div>');

    // 将留言添加到div元素
    div.html(`
    <span><a href="${comment.website}">${comment.username}</a> 说: <br> <span class="context_remsg"> ${comment.content}</span></p>
    <button class="reply-button">Reply</button>
    <div class="reply-form"></div>
    <div class="replies"></div>
  `);

    // 当用户点击回复按钮时，显示一个回复表单
    div.find('.reply-button').on('click', function() {
        const replyForm = div.find('.reply-form:first');
         replyForm.append(commentForm);
        parentIdInput.val(comment.id);
        rootIdInput.val(comment.rootId || comment.id);
    });

    // 返回新的div元素
    return div;
}


function loadComments() {
    const page = getPageFromUrl();
    // 向后端API发送一个GET请求，获取所有的留言
    $.ajax({
        url: '/comments',
        method: 'GET',
        contentType: 'application/json',
        data: {
            page: page,
        },
        success: function(comments) {
            // 清空留言列表
            commentList.empty();

            // 创建一个映射，将每个留言的ID映射到它的HTML元素
            const commentElements = {};
            for (const comment of comments) {
                commentElements[comment.id] = createCommentElement(comment);
            }

            // 将每条留言添加到它的父留言的回复列表，或者如果它没有父留言，就添加到留言列表
            for (const comment of comments) {
                const parentElement = comment.parentId ? commentElements[comment.parentId].find('.replies:first') : commentList;
                parentElement.append(commentElements[comment.id]);
            }
        }
    });
}

// 页面加载时，加载所有的留言
loadComments();
